package cmd

import (
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

var RequirementsB64 string = ""

func GetRequirements() string {
    if RequirementsB64 == "" {
        return ""
    }

    decoded, err := base64.StdEncoding.DecodeString(RequirementsB64)
    if err != nil {
        return ""
    }

    return strings.TrimSpace(string(decoded))
}

type Requirement struct {
    Packages   []string
    MinVersion string
    MaxVersion string
}

func parseRequirements(content string) []Requirement {
    var reqs []Requirement
    lines := strings.Split(content, "\n")

    for _, line := range lines {
        line = strings.TrimSpace(line)

        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        alternatives := strings.Split(line, "|")
        var packages []string
        var minVer, maxVer string

        for i, alt := range alternatives {
            alt = strings.TrimSpace(alt)
            parts := strings.Fields(alt)

            if len(parts) < 3 {
                continue
            }

            packages = append(packages, parts[0])

            if i == 0 {
                minVer = parts[1]
                maxVer = parts[2]
            }
        }

        if len(packages) > 0 {
            reqs = append(reqs, Requirement{
                Packages:   packages,
                MinVersion: minVer,
                MaxVersion: maxVer,
            })
        }
    }

    return reqs
}

func checkPackageExists(pkg string) bool {
    _, err := exec.LookPath(pkg)
    return err == nil
}

func formatVersionInfo(minVer, maxVer string) string {
    hasMin := minVer != "" && minVer != "*"
    hasMax := maxVer != "" && maxVer != "*"

    if !hasMin && !hasMax {
        return ""
    }

    if hasMin && hasMax {
        return fmt.Sprintf("(recommended: %s - %s)", minVer, maxVer)
    }

    if hasMin {
        return fmt.Sprintf("(recommended: %s+)", minVer)
    }

    return fmt.Sprintf("(recommended: up to %s)", maxVer)
}

func CheckRequirements() {
    content := GetRequirements()
    if content == "" {
        Print("No requirements information available.")
        return
    }

    reqs := parseRequirements(content)
    if len(reqs) == 0 {
        Print("No requirements to check.")
        return
    }

    Print("KVMage System Requirements Check")
    Print("=================================")
    Print("")

    allMet := true
    warnings := []string{}

    for _, req := range reqs {
        versionInfo := formatVersionInfo(req.MinVersion, req.MaxVersion)

        if len(req.Packages) == 1 {
            pkg := req.Packages[0]
            if checkPackageExists(pkg) {
                if versionInfo != "" {
                    Print("  ✓ %s %s", pkg, versionInfo)
                } else {
                    Print("  ✓ %s", pkg)
                }
            } else {
                if versionInfo != "" {
                    Print("  ✗ %s (not found) %s", pkg, versionInfo)
                } else {
                    Print("  ✗ %s (not found)", pkg)
                }
                allMet = false
            }
        } else {
            found := false
            var foundPkg string
            for _, pkg := range req.Packages {
                if checkPackageExists(pkg) {
                    found = true
                    foundPkg = pkg
                    break
                }
            }

            if found {
                Print("  ✓ %s (alternatives: %s)", foundPkg, strings.Join(req.Packages, ", "))
            } else {
                Print("  ✗ Missing all alternatives: %s", strings.Join(req.Packages, ", "))
                allMet = false
            }
        }
    }

    Print("")
    if allMet {
        Print("Result: All requirements met ✓")
    } else {
        Print("Result: Missing required dependencies ✗")
        os.Exit(1)
    }

    if len(warnings) > 0 {
        Print("")
        Print("Warnings:")
        for _, w := range warnings {
            Print("  - %s", w)
        }
    }

    Print("")
    Print("Note: Listed versions are tested and supported. Other versions are unsupported and used at your own risk.")
}

func QuickCheckRequirements() error {
    content := GetRequirements()
    if content == "" {
        return nil
    }

    reqs := parseRequirements(content)
    var missing []string

    for _, req := range reqs {
        if len(req.Packages) == 1 {
            if !checkPackageExists(req.Packages[0]) {
                missing = append(missing, req.Packages[0])
            }
        } else {
            found := false
            for _, pkg := range req.Packages {
                if checkPackageExists(pkg) {
                    found = true
                    break
                }
            }
            if !found {
                missing = append(missing, fmt.Sprintf("(%s)", strings.Join(req.Packages, " or ")))
            }
        }
    }

    if len(missing) > 0 {
        return fmt.Errorf("missing required dependencies: %s", strings.Join(missing, ", "))
    }

    return nil
}
