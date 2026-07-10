package ports

import (
	"context"
	"io"

	"servermanager/internal/domain"
)

// ImageBuilder turns checked-out source into a locally tagged image. Build
// output (generator + daemon build stream) is written to logSink as it
// arrives; logSink must be safe for concurrent writes. A nil error means an
// image tagged in.ImageRef exists in the daemon afterwards.
type ImageBuilder interface {
	Build(ctx context.Context, in domain.BuildInput, logSink io.Writer) error
}
