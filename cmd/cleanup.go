package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CleanupArtifacts() {
	paths := []string{
		TempImagePath,
		TempInstallFile,
		TempInstallMedia,
		TempImageSource,
	}
 
	for _, path := range paths {
		if path == "" {
			continue
		}

		if !isSafeTempPath(path) {
			PrintVerbose(2, "Skipping unsafe path: %s", path)
			continue
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		if err := os.Remove(path); err != nil {
			PrintError("Warning: failed to remove temp file %s: %v", path, err)
		} else {
			PrintVerbose(2, "Removed temporary file: %s", path)
		}
	}

	if TempUploadDir != "" {
		if _, err := os.Stat(TempUploadDir); err == nil {
			if err := os.RemoveAll(TempUploadDir); err != nil {
				PrintError("Warning: failed to remove upload temp directory %s: %v", TempUploadDir, err)
			} else {
				PrintVerbose(2, "Removed upload temp directory: %s", TempUploadDir)
			}
		}
	} 

	if TempImageName != "" {
		removeTempVM(TempImageName)
	}
}

func isSafeTempPath(path string) bool {
	return strings.HasPrefix(path, "/var/lib/libvirt/images/kvmage-")
}

func removeTempVM(vmName string) {
	PrintVerbose(2, "Checking for temporary VM: %s", vmName)

	checkCmd := exec.Command("virsh", "dominfo", vmName)
	if err := checkCmd.Run(); err != nil {
		PrintVerbose(2, "Domain %s does not exist. Skipping undefine.", vmName)
		return
	}

	PrintVerbose(2, "Undefining temporary VM: %s", vmName)

	cmd := exec.Command("virsh", "undefine", "--nvram", vmName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		PrintError("Warning: failed to undefine VM %s: %v", vmName, err)
	}
}

func CleanupOrphanedTempFiles() {
	var toDeleteFiles []string
	var toDeleteDirs []string
    var orphanedVMs []string
	var totalSize int64

	imagesDir := "/var/lib/libvirt/images"
    if entries, err := os.ReadDir(imagesDir); err == nil {
    	for _, entry := range entries {
    		name := entry.Name()
    		fullPath := filepath.Join(imagesDir, name)


			if !strings.HasPrefix(name, "kvmage-") {
				continue
			}
			
			info, err := entry.Info()
			if err != nil {
				continue
			}		

			if info.IsDir() {
				dirSize := getDirSize(fullPath)
				toDeleteDirs = append(toDeleteDirs, fullPath)
				totalSize += dirSize
			} else {
				if isFileInUse(fullPath) {
				PrintVerbose(2, "Skipping in-use file: %s", fullPath)
					continue
				}
				toDeleteFiles = append(toDeleteFiles, fullPath)
				totalSize += info.Size()
			}
		}
	} else {
		PrintError("Failed to read directory %s: %v", imagesDir, err)
	}
		 
	tmpDir := "/var/lib/libvirt/tmp"
	if entries, err := os.ReadDir(tmpDir); err == nil {
		for _, entry := range entries {
			name := entry.Name()
			fullPath := filepath.Join(tmpDir, name)
	
			if !strings.HasPrefix(name, "kvmage-") {
				continue
			}
	
			info, err := entry.Info()
			if err != nil {
				continue
			}
	
			if info.IsDir() {
				dirSize := getDirSize(fullPath)
				toDeleteDirs = append(toDeleteDirs, fullPath)
				totalSize += dirSize
			} else {
				if isFileInUse(fullPath) {
					PrintVerbose(2, "Skipping in-use file: %s", fullPath)
					continue
				}
				toDeleteFiles = append(toDeleteFiles, fullPath)
				totalSize += info.Size()
			}
		} 
	}

	orphanedVMs = findOrphanedVMs()
	
	totalItems := len(toDeleteFiles) + len(toDeleteDirs) + len(orphanedVMs)
	if totalItems == 0 {
		Print("No orphaned kvmage artifacts found.")
		return
	}

	Print("Found orphaned kvmage artifacts:\n")
	
	if len(toDeleteFiles) > 0 {
		Print("\nFiles (%d):", len(toDeleteFiles))
		for _, path := range toDeleteFiles {
			size := getFileSize(path)
			Print("  %-65s %8s", path, formatSize(size))
		} 
	}
	
	if len(toDeleteDirs) > 0 {
		Print("\nDirectories (%d):", len(toDeleteDirs))
		for _, path := range toDeleteDirs {
			size := getDirSize(path)
			Print("  %-65s %8s", path, formatSize(size))
		}
	}
	
	if len(orphanedVMs) > 0 {
		Print("\nOrphaned VMs (%d):", len(orphanedVMs))
		for _, vm := range orphanedVMs {
			Print("  %s", vm)
		}
	}

	Print("\nTotal reclaimable space: %s", formatSize(totalSize))

	Print("\nDo you want to delete these artifacts? [y/N]: ")
	var input string
	fmt.Scanln(&input)

	if strings.ToLower(input) != "y" {
		Print("Aborted. No artifacts deleted.")
		return
	}

	for _, path := range toDeleteFiles {
		 if err := os.Remove(path); err != nil {
			PrintError("Failed to delete %s: %v", path, err)
		} else {
			Print("Deleted: %s", path)
		}
	}
	
	for _, path := range toDeleteDirs {
		if err := os.RemoveAll(path); err != nil {
			PrintError("Failed to delete directory %s: %v", path, err)
		} else {
			Print("Deleted directory: %s", path)
		}
	}
	
	for _, vm := range orphanedVMs {
		removeTempVM(vm)
		Print("Removed VM: %s", vm)
	}
}

func findOrphanedVMs() []string {
	var orphaned []string
	
	cmd := exec.Command("virsh", "list", "--all", "--name")
	output, err := cmd.Output()
	if err != nil {
		PrintVerbose(2, "Failed to list VMs: %v", err)
		return orphaned
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, name := range lines {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if strings.HasPrefix(name, "kvmage-") {
			checkCmd := exec.Command("virsh", "domstate", name)
			state, _ := checkCmd.Output()
			if strings.TrimSpace(string(state)) != "running" {
				orphaned = append(orphaned, name)
			} else {
				PrintVerbose(2, "Skipping running VM: %s", name)
			}
		}
	}
	
		return orphaned
	}
	
	func getDirSize(path string) int64 {
		var size int64
		filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				size += info.Size()
			}
		return nil
	})
	return size
}

func isFileInUse(path string) bool {
	cmd := exec.Command("lsof", "-Fn", "--", path)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func formatSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
