package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func RunCustomize(opts *Options, tempName, tempPath string) error {
	RequireRoot()

	PrintVerbose(1, "Starting customize mode...")
	PrintVerbose(2, "Customizing image at path: %s", tempPath)

	if opts.ImageSize != "" {
		PrintVerbose(1, "Checking image size...")

		cmd := exec.Command("qemu-img", "info", "--output=json", tempPath)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get image info: %w", err)
		}

		var info struct {
			VirtualSize int64 `json:"virtual-size"`
		}
		if err := json.Unmarshal(output, &info); err != nil {
			return fmt.Errorf("failed to parse qemu-img info output: %w", err)
		}
		currentSizeGB := float64(info.VirtualSize) / (1024 * 1024 * 1024)

		targetSizeStr := strings.TrimSuffix(opts.ImageSize, "G")
		targetSizeGB, err := strconv.ParseFloat(targetSizeStr, 64)
		if err != nil {
			return fmt.Errorf("invalid image_size format: %w", err)
		}

		if targetSizeGB < currentSizeGB {
			PrintError(
				"Refusing to shrink image: requested size %.1f GB is less than current size %.1f GB",
				targetSizeGB,
				currentSizeGB,
			)
			return fmt.Errorf("image_size must not be smaller than the current image size")
		} else if targetSizeGB > currentSizeGB {
			PrintVerbose(
				1,
				"Expanding image from %.1f GB to %.1f GB",
				currentSizeGB,
				targetSizeGB,
			)
			if err := exec.Command("qemu-img", "resize", tempPath, opts.ImageSize).Run(); err != nil {
				return fmt.Errorf("qemu-img resize failed: %w", err)
			}
		} else {
			PrintVerbose(2, "Image is already at requested size: %.1f GB", currentSizeGB)
		}
	}

	if opts.ImagePartition != "" {
		PrintVerbose(1, "Expanding partition: %s", opts.ImagePartition)

		backupPath := tempPath + ".orig"
		if err := os.Rename(tempPath, backupPath); err != nil {
			return fmt.Errorf("failed to rename original image for virt-resize: %w", err)
		}

		args := []string{
			"--expand", opts.ImagePartition,
			backupPath,
			tempPath,
		}

		PrintVerbose(
			2,
			"Running virt-resize with args: virt-resize %s",
			joinArgs(args),
		)
		if err := exec.Command("virt-resize", args...).Run(); err != nil {
			return fmt.Errorf("virt-resize failed: %w", err)
		}
	}

	args := []string{"-a", tempPath}
	if verboseLevel >= 1 {
		args = append(args, "-v")
	}
	if verboseLevel >= 2 {
		args = append(args, "-x")
	}
	for i, originalPath := range opts.Upload {
		tempUploadPath := TempUploadPaths[i]
		vmPath := filepath.Join("/tmp/kvmage", originalPath)

		info, err := os.Stat(tempUploadPath)
		if err != nil {
			return fmt.Errorf("failed to stat upload path %s: %w", originalPath, err)
		}
		if info.IsDir() {
			vmParent := filepath.Dir(vmPath)
			args = append(args, "--run-command", fmt.Sprintf("mkdir -p %s", vmParent))
			args = append(args, "--copy-in", fmt.Sprintf("%s:%s", tempUploadPath, vmParent))
			PrintVerbose(2, "Uploading directory: %s -> %s", originalPath, vmPath)
		} else {
			vmParentDir := filepath.Dir(vmPath)
			args = append(args, "--run-command", fmt.Sprintf("mkdir -p %s", vmParentDir))
			args = append(args, "--upload", fmt.Sprintf("%s:%s", tempUploadPath, vmPath))
			PrintVerbose(2, "Uploading file: %s -> %s", originalPath, vmPath)
		}
	}
	for _, execPath := range opts.Execute {
		vmPath := filepath.Join("/tmp/kvmage", execPath)
		args = append(args, "--chmod", fmt.Sprintf("0755:%s", vmPath))
		PrintVerbose(2, "Setting executable permissions for: %s", vmPath)
		args = append(args, "--run-command", vmPath)
		PrintVerbose(2, "Will execute: %s", vmPath)
	}
	if len(opts.Upload) > 0 {
		args = append(args, "--run-command", "rm -rf /tmp/kvmage || true")
		PrintVerbose(2, "Will cleanup /tmp/kvmage in VM")
	}
	if opts.Hostname != "" {
		args = append(args, "--hostname", opts.Hostname)
		PrintVerbose(2, "Setting hostname: %s", opts.Hostname)
	}

	PrintVerbose(
		3,
		"Running virt-customize with args: virt-customize %s",
		joinArgs(args),
	)
	PrintVerbose(1, "Executing virt-customize...")

	run := func(cmdName string, args []string) error {
		cmd := exec.Command(cmdName, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if err := run("virt-customize", args); err != nil {
		return fmt.Errorf("virt-customize failed: %w", err)
	}

	return nil
}
