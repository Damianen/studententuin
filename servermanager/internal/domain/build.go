package domain

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
