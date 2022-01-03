package horcrux

import (
	"bufio"
	"os/exec"
	"strings"
)

func smartcardHasKey(keyID string) bool {
	if !smartcardIsAttached() {
		return false
	}
	cmd := exec.Command("gpg", "--with-colons", "--card-status")
	stdout, stderr := cmd.Output()
	if stderr != nil {
		return false
	}

	smartcardFingerprints := parseSmartcardFingerprints(stdout)
	for _, fp := range smartcardFingerprints {
		if fp[len(fp)-16:] == keyID {
			return true
		}
	}
	return false
}

func smartcardIsAttached() bool {
	cmd := exec.Command("gpg", "--card-status")
	_, stderr := cmd.Output()
	return stderr == nil
}

func parseSmartcardFingerprints(rawStdout []byte) []string {
	stdout := string(rawStdout)
	smartcardFingerprints := []string{}

	sc := bufio.NewScanner(strings.NewReader(stdout))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "fpr:") {
			components := strings.Split(line, ":")
			for i, fp := range components {
				// Fingerprints are 40 chars
				if len(fp) == 40 {
					// Only exist in positoins 1, 2 or 3
					// ref: https://helpful.wiki/gpg
					if i == 1 || i == 2 || i == 3 {
						smartcardFingerprints = append(smartcardFingerprints, fp)
					}
				}
			}
		}
	}
	return smartcardFingerprints
}
