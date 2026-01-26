package cmd

import (
        "errors"
        "os"
)

// DefaultConfigFiles are the filenames kvmage looks for when -f is used without a value
var DefaultConfigFiles = []string{
        "kvmage.yml",
        "kvmage.yaml",
}

// FindDefaultConfig looks for a default config file in the current directory
func FindDefaultConfig() string {
        for _, filename := range DefaultConfigFiles {
                if _, err := os.Stat(filename); err == nil {
                        return filename
                }
        }
        return ""
}

// ResolveConfigPath resolves "AUTO" to a default config file, or returns the path as-is
func ResolveConfigPath(path string) (string, error) {
        if path != "AUTO" {
                return path, nil
        }
        if defaultConfig := FindDefaultConfig(); defaultConfig != "" {
                return defaultConfig, nil
        }
        return "", errors.New("no config file specified and no kvmage.yml or kvmage.yaml found in current directory")
}

func ValidateModeFlags(runMode bool, configPath string) error {
        switch {
        case runMode && configPath != "":
                return errors.New("cannot specify both --run and --config")
        case !runMode && configPath == "":
                return errors.New("must specify either --run or --config")
        default:
                return nil
        }
}