package cmd

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	preseedRegex     = regexp.MustCompile(`^d-i\b`)
	ksSectionRegex   = regexp.MustCompile(`^%(packages|pre|post|addon|include)\b`)
	ksDirectiveRegex = regexp.MustCompile(`^(auth|authconfig|autopart|bootloader|cdrom|clearpart|cmdline|device|deviceprobe|firewall|firstboot|ignoredisk|install|keyboard|lang|logging|logvol|mediacheck|network|part|partdump|partinfo|poweroff|reboot|rootpw|selinux|services|shutdown|skipx|text|timezone|url|user|volgroup|xconfig)\b`)
)

func DetectInstallType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	preseedCount := 0
	ksSectionCount := 0
	ksDirectiveCount := 0

	const maxLines = 300

	for i := 0; i < maxLines && scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, "//") ||
			strings.HasPrefix(line, ";") {
			continue
		}

		if preseedRegex.MatchString(line) {
			preseedCount++
		}

		if ksSectionRegex.MatchString(line) {
			ksSectionCount++
		}

		if ksDirectiveRegex.MatchString(line) {
			ksDirectiveCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	preseedScore := preseedCount
	kickstartScore := ksSectionCount*5 + ksDirectiveCount

	PrintVerbose(3, "DetectInstallType(%s): preseedScore=%d kickstartScore=%d (preseedCount=%d ksSection=%d ksDirective=%d)",
		path, preseedScore, kickstartScore, preseedCount, ksSectionCount, ksDirectiveCount)

	scores := map[string]int{
		"preseed":   preseedScore,
		"kickstart": kickstartScore,
	}

	bestType := "unknown"
	bestScore := 0
	secondScore := 0

	for t, s := range scores {
		if s > bestScore {
			secondScore = bestScore
			bestScore = s
			bestType = t
		} else if s > secondScore {
			secondScore = s
		}
	}

	if bestScore == 0 || bestScore == secondScore {
		PrintVerbose(2, "DetectInstallType(%s): unknown install type", path)
		return "unknown", nil
	}

	PrintVerbose(2, "Detected install type for %s: %s (score=%d)", path, bestType, bestScore)
	return bestType, nil
}
