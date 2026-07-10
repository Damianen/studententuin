package domain

import "fmt"

// BuildInput is everything the ImageBuilder needs to turn a checkout into a
// locally tagged image. Dir and ImageRef are manager-generated; Build- and
// StartCommand are shell strings by design — they execute inside build/user
// containers, never on the manager host (only /run takes an argv array).
type BuildInput struct {
	Dir          string
	ImageRef     string
	Env          map[string]string
	BuildCommand string
	StartCommand string
	Port         int
}

// ValidateCommand rejects control characters in a build/start command —
// they are single-line shell strings that end up in a generated Dockerfile.
func ValidateCommand(field, cmd string) error {
	for _, r := range cmd {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("%s contains control characters: %w", field, ErrInvalid)
		}
	}
	return nil
}
