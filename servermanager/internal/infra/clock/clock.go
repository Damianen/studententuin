// Package clock provides the real ports.Clock implementation.
package clock

import "time"

// System reads the wall clock.
type System struct{}

func (System) Now() time.Time { return time.Now() }
