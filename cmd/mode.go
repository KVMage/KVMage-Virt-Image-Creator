package cmd

import (
    "errors"
    "os"
)

var DefaultConfigFiles = []string{
    "kvmage.yml",
    "kvmage.yaml",
}

func FindDefaultConfig() string {
    for _, filename := range DefaultConfigFiles {
        if _, err := os.Stat(filename); err == nil {
            return filename
        }
    }
    return ""
}

func ResolveConfigPath(path string, args []string) (string, error) {
    if path == "AUTO" && len(args) > 0 {
        return args[0], nil
    }
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