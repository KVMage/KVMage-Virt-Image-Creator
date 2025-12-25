package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	TempImageName    string
	TempImagePath    string
	TempInstallFile  string
	TempInstallMedia string
	TempImageSource  string
	TempUploadPaths  []string
	TempUploadDir    string
)

const tempDir = "/var/lib/libvirt/images"
const uploadTempDir = "/var/lib/libvirt/tmp"

func CreateTempImage(opts *Options) (string, string, error) {
	RequireRoot()

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate random name: %w", err)
	}
	TempImageName = fmt.Sprintf("kvmage-%s", hex.EncodeToString(b))
	TempImagePath = filepath.Join(tempDir, TempImageName+".qcow2")

	switch opts.VirtMode {

	case "install":
		PrintVerbose(2, "Temporary image path: %s", TempImagePath)
		PrintVerbose(2, "Requested image size: %s", opts.ImageSize)

		args := []string{"create", "-f", "qcow2", "-o", "compat=0.10", TempImagePath, opts.ImageSize}
		PrintVerbose(3, "Running command: qemu-img %s", joinArgs(args))

		if err := exec.Command("qemu-img", args...).Run(); err != nil {
			return "", "", fmt.Errorf("qemu-img create failed: %w", err)
		}

	case "customize":
		PrintVerbose(2, "Copying source image %s to %s", opts.ImageSource, TempImagePath)

		args := []string{"convert", "-O", "qcow2", opts.ImageSource, TempImagePath}
		PrintVerbose(3, "Running command: qemu-img %s", joinArgs(args))

		if err := exec.Command("qemu-img", args...).Run(); err != nil {
			return "", "", fmt.Errorf("qemu-img convert failed: %w", err)
		}

	default:
		return "", "", fmt.Errorf("invalid VirtMode: must be 'install' or 'customize'")
	}

	return TempImageName, TempImagePath, nil
}

func resolveInstallMedia(src string) (string, error) {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		if strings.HasSuffix(strings.ToLower(src), ".iso") {
			filename := filepath.Base(src)
			dest := filepath.Join(os.TempDir(), filename)

			if _, err := os.Stat(dest); err == nil {
				PrintVerbose(2, "Using cached remote ISO: %s", dest)
				return dest, nil
			}

			var downloader string
			if _, err := exec.LookPath("curl"); err == nil {
				downloader = "curl"
			} else if _, err := exec.LookPath("wget"); err == nil {
				downloader = "wget"
			} else {
				return "", fmt.Errorf("neither curl nor wget is installed")
			}

			PrintVerbose(2, "Downloading ISO from %s using %s", src, downloader)

			var cmd *exec.Cmd
			if downloader == "curl" {
				cmd = exec.Command("curl", "-L", "-v", "-o", dest, src)
			} else {
				cmd = exec.Command("wget", "-v", "-O", dest, src)
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return "", fmt.Errorf("failed to download remote ISO: %w", err)
			}

			return dest, nil
		} else {
			return src, nil
		}
	}

	info, err := os.Stat(src)
	if err != nil {
		return "", fmt.Errorf("cannot access install media: %w", err)
	}

	if info.IsDir() {
		absPath, err := filepath.Abs(src)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	return src, nil
}

func copyToTemp(src string, label string) (string, error) {
	if src == "" {
		return "", nil
	}

	ext := filepath.Ext(src)
	destName := fmt.Sprintf("%s-%s.temp%s", TempImageName, label, ext)
	dest := filepath.Join(tempDir, destName)

	PrintVerbose(2, "Copying %s to temp file: %s", label, dest)

	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file %s: %w", dest, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return "", fmt.Errorf("failed to copy to temp file %s: %w", dest, err)
	}

	return dest, nil
}

func CopyInputFilesToTempDir(opts *Options) error {
	RequireRoot()

	var err error

	if len(opts.Upload) > 0 {
		TempUploadDir = filepath.Join(uploadTempDir, TempImageName)
		if err := os.MkdirAll(TempUploadDir, 0755); err != nil {
			return fmt.Errorf("failed to create upload temp directory: %w", err)
		}
		PrintVerbose(2, "Created upload temp directory: %s", TempUploadDir)
	}
	 
	if TempInstallFile, err = copyToTemp(opts.InstallFile, "auto"); err != nil {
		return fmt.Errorf("auto install file copy failed: %w", err)
	}

	if opts.InstallMedia != "" {
		var resolved string
		if resolved, err = resolveInstallMedia(opts.InstallMedia); err != nil {
			return fmt.Errorf("install media resolution failed: %w", err)
		}

		if strings.HasPrefix(resolved, "http://") || strings.HasPrefix(resolved, "https://") {
			TempInstallMedia = resolved
			PrintVerbose(2, "Using remote install media: %s", TempInstallMedia)
		} else if info, err := os.Stat(resolved); err == nil && info.IsDir() {
			TempInstallMedia = resolved
			PrintVerbose(2, "Using local install tree: %s", TempInstallMedia)
		} else {
			if TempInstallMedia, err = copyToTemp(resolved, "iso"); err != nil {
				return fmt.Errorf("install media copy failed: %w", err)
			}
		}
	}

	if TempImageSource, err = copyToTemp(opts.ImageSource, "src"); err != nil {
		return fmt.Errorf("source image copy failed: %w", err)
	}
	for i, uploadPath := range opts.Upload {
		label := fmt.Sprintf("upload-%d", i)
		tempPath, err := copyToTempCustomDir(uploadPath, label, TempUploadDir)
		if err != nil {
			return fmt.Errorf("upload file copy failed for %s: %w", uploadPath, err)
		}
		TempUploadPaths = append(TempUploadPaths, tempPath)
	}

	return nil
}

func copyToTempCustomDir(src string, label string, destDir string) (string, error) {
	if src == "" {
		return "", nil
	}

	info, err := os.Stat(src)
	if err != nil {
		return "", fmt.Errorf("failed to stat %s: %w", src, err)
	}

	if info.IsDir() {
		destName := fmt.Sprintf("%s-%s.temp", TempImageName, label)
		dest := filepath.Join(destDir, destName)
	
		PrintVerbose(2, "Copying %s directory to temp: %s", label, dest)
	
		if err := copyDir(src, dest); err != nil {
			return "", fmt.Errorf("failed to copy directory %s: %w", src, err)
		}
		return dest, nil
	}

	func copyDir(src string, dst string) error {
		srcInfo, err := os.Stat(src)
		if err != nil {
			return err
		}
	
		if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
			return err
		}
	
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
	
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())
	
			if entry.IsDir() {
				if err := copyDir(srcPath, dstPath); err != nil {
					return err
				}
			} else {
				if err := copyFile(srcPath, dstPath); err != nil {
					return err
				}
			}
		}
	
		return nil
	}
	
	func copyFile(src string, dst string) error {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()
	
		srcInfo, err := srcFile.Stat()
		if err != nil {
			return err
		}
	
		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()
	
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	
		return os.Chmod(dst, srcInfo.Mode())
	}
 

	ext := filepath.Ext(src)
	destName := fmt.Sprintf("%s-%s.temp%s", TempImageName, label, ext)
	dest := filepath.Join(destDir, destName)

	PrintVerbose(2, "Copying %s to temp file: %s", label, dest)

	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file %s: %w", dest, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return "", fmt.Errorf("failed to copy to temp file %s: %w", dest, err)
	}

	return dest, nil
}
 
func joinArgs(args []string) string {
	return fmt.Sprintf("%q", args)
}
