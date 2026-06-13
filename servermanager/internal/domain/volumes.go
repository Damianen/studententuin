package domain

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// volumeNameRe matches Docker's volume-name charset. A name cannot start with
// '/' or '.', which structurally rules out host bind mounts (§3.3).
var volumeNameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

// ParseVolumeSpec validates a "name:/abs/target" volume entry and returns the
// named-volume name and the absolute container path it mounts to.
func ParseVolumeSpec(entry string) (name, target string, err error) {
	name, target, ok := strings.Cut(entry, ":")
	if !ok {
		return "", "", fmt.Errorf("volume %q: want \"name:/abs/target\"", entry)
	}
	if !volumeNameRe.MatchString(name) {
		return "", "", fmt.Errorf("volume %q: invalid volume name", entry)
	}
	if !strings.HasPrefix(target, "/") {
		return "", "", fmt.Errorf("volume %q: target must be an absolute path", entry)
	}
	if strings.Contains(target, "..") || path.Clean(target) != target {
		return "", "", fmt.Errorf("volume %q: target must be a clean path", entry)
	}
	return name, target, nil
}
