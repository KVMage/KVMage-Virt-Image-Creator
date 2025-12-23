package cmd

import (
    "os"

    "github.com/spf13/cobra"
)

func runMain(cmd *cobra.Command, args []string) error {
    if HandleGlobalFlags() {
        return nil
    }

    if err := QuickCheckRequirements(); err != nil {
        PrintError("WARNING: %v", err)
        PrintError("Run 'kvmage --check-requirements' for details")
        PrintError("")
    }

    if err := ValidateModeFlags(runMode, configPath); err != nil {
        PrintError("%v", err)
        os.Exit(1)
    }

    SetupSignalHandler(CleanupArtifacts)
    defer CleanupArtifacts()

    if runMode {
        return runRunMode()
    } else if configPath != "" {
        return runConfigMode(configPath)
    }

    return nil
}
