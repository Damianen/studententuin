package domain

import "testing"

func TestParseVolumeSpec(t *testing.T) {
	valid := []struct {
		in         string
		wantName   string
		wantTarget string
	}{
		{"data:/var/lib/data", "data", "/var/lib/data"},
		{"app_cache:/tmp/cache", "app_cache", "/tmp/cache"},
		{"v1.backup:/backup", "v1.backup", "/backup"},
		{"0abc-def:/srv", "0abc-def", "/srv"},
	}
	for _, tc := range valid {
		name, target, err := ParseVolumeSpec(tc.in)
		if err != nil {
			t.Errorf("ParseVolumeSpec(%q) returned error: %v", tc.in, err)
			continue
		}
		if name != tc.wantName || target != tc.wantTarget {
			t.Errorf("ParseVolumeSpec(%q) = (%q, %q), want (%q, %q)", tc.in, name, target, tc.wantName, tc.wantTarget)
		}
	}

	invalid := []string{
		"",                       // empty
		"data",                   // no colon
		":/var/lib/data",         // empty name
		"data:",                  // empty target
		"/host/path:/data",       // bind mount: name starts with /
		"./rel:/data",            // bind mount: name starts with .
		".hidden:/data",          // name starts with .
		"bad name:/data",         // space in name
		"data:relative/path",     // relative target
		"data:/var/../etc",       // traversal
		"data:/var//lib",         // not clean
		"data:/var/lib/",         // trailing slash, not clean
		"da$ta:/data",            // invalid char
		"data:/var/lib/..",       // traversal at end
	}
	for _, in := range invalid {
		if name, target, err := ParseVolumeSpec(in); err == nil {
			t.Errorf("ParseVolumeSpec(%q) = (%q, %q), want error", in, name, target)
		}
	}
}
