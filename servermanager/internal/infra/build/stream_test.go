package build

import (
	"io"
	"strings"
	"testing"
)

func TestDecodeBuildStreamSuccess(t *testing.T) {
	body := io.NopCloser(strings.NewReader(
		`{"stream":"Step 1/4 : FROM base\n"}` + "\n" +
			`{"status":"Pulling from railwayapp/nixpacks","id":"ubuntu"}` + "\n" +
			`{"status":"Downloading","id":"layer1","progressDetail":{"current":10,"total":100}}` + "\n" +
			`{"stream":"Successfully built abc123\n"}` + "\n",
	))

	var sink strings.Builder
	if err := decodeBuildStream(body, &sink); err != nil {
		t.Fatalf("decodeBuildStream: %v", err)
	}

	out := sink.String()
	for _, want := range []string{"Step 1/4", "Successfully built abc123", "ubuntu: Pulling from railwayapp/nixpacks"} {
		if !strings.Contains(out, want) {
			t.Errorf("sink missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Downloading") {
		t.Errorf("progress spam reached the sink:\n%s", out)
	}
}

func TestDecodeBuildStreamError(t *testing.T) {
	body := io.NopCloser(strings.NewReader(
		`{"stream":"Step 3/4 : RUN npm i\n"}` + "\n" +
			`{"errorDetail":{"code":1,"message":"The command 'npm i' returned a non-zero code: 1"},"error":"..."}` + "\n" +
			`{"stream":"tail after error\n"}` + "\n",
	))

	var sink strings.Builder
	err := decodeBuildStream(body, &sink)
	if err == nil || !strings.Contains(err.Error(), "non-zero code: 1") {
		t.Fatalf("err = %v, want the daemon's error message", err)
	}
	// The stream is drained even after the error message arrives.
	if !strings.Contains(sink.String(), "tail after error") {
		t.Errorf("stream tail lost after error:\n%s", sink.String())
	}
}

func TestDecodeBuildStreamGarbage(t *testing.T) {
	body := io.NopCloser(strings.NewReader(`{"stream":"ok\n"}` + "\n" + `not json at all`))
	if err := decodeBuildStream(body, io.Discard); err == nil || !strings.Contains(err.Error(), "build stream") {
		t.Errorf("err = %v, want a decode error", err)
	}
}
