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
    Packages []string
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

        for _, alt := range alternatives {
            alt = strings.TrimSpace(alt)
        }

        if len(packages) > 0 {
            reqs = append(reqs, Requirement{
                Packages: packages,
            })
        }
    }

    return reqs
}

func checkPackageExists(pkg string) bool {
    _, err := exec.LookPath(pkg)
    return err == nil
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

    for _, req := range reqs {
        if len(req.Packages) == 1 {
            pkg := req.Packages[0]
            if checkPackageExists(pkg) {
                Print("  ✓ %s", pkg)
            } else {
                Print("  ✗ %s (not found)", pkg)
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
