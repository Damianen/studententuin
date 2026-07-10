package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/moby/moby/api/types/jsonstream"
)

// decodeBuildStream copies the daemon's build output into sink and surfaces
// the terminal error, if any. A clean EOF without an errorDetail message is
// success — that is the legacy builder's contract.
func decodeBuildStream(body io.ReadCloser, sink io.Writer) error {
	defer func() { _ = body.Close() }()

	dec := json.NewDecoder(body)
	var buildErr error
	for {
		var msg jsonstream.Message
		if err := dec.Decode(&msg); err != nil {
			if errors.Is(err, io.EOF) {
				return buildErr
			}
			if buildErr != nil {
				return buildErr
			}
			return fmt.Errorf("reading build stream: %w", err)
		}

		if msg.Error != nil && buildErr == nil {
			// Remember and keep draining so the tail of the log still lands
			// in the sink.
			buildErr = fmt.Errorf("build failed: %s", msg.Error.Message)
		}
		switch {
		case msg.Stream != "":
			_, _ = io.WriteString(sink, msg.Stream)
		case msg.Status != "" && msg.Progress == nil:
			// Progress ticks (layer download counters) are spam; the
			// status transitions are worth keeping.
			if msg.ID != "" {
				_, _ = fmt.Fprintf(sink, "%s: %s\n", msg.ID, msg.Status)
			} else {
				_, _ = fmt.Fprintln(sink, msg.Status)
			}
		}
	}
}
