package build

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteTar(t *testing.T) {
	dir := t.TempDir()
	mustWrite := func(name string, mode os.FileMode) {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("content of "+name), mode); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("package.json", 0o644)
	mustWrite(".nixpacks/Dockerfile", 0o644)
	mustWrite(".nixpacks/assets/run.sh", 0o755)
	mustWrite(".git/config", 0o644) // must be skipped
	mustWrite("secret.txt", 0o644)
	if err := os.Symlink("secret.txt", filepath.Join(dir, "link.txt")); err != nil {
		t.Fatal(err)
	}
	// A symlink pointing outside the context must be emitted as a link, not
	// followed (that would exfiltrate host files into the image).
	if err := os.Symlink("/etc/hostname", filepath.Join(dir, "escape.txt")); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := writeTar(&buf, dir); err != nil {
		t.Fatalf("writeTar: %v", err)
	}

	entries := map[string]*tar.Header{}
	bodies := map[string]string{}
	tr := tar.NewReader(&buf)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("reading tar: %v", err)
		}
		entries[hdr.Name] = hdr
		if hdr.Typeflag == tar.TypeReg {
			body, _ := io.ReadAll(tr)
			bodies[hdr.Name] = string(body)
		}
	}

	for _, name := range []string{"package.json", ".nixpacks/Dockerfile", ".nixpacks/assets/run.sh", "secret.txt"} {
		if entries[name] == nil {
			t.Errorf("%s missing from context", name)
		}
	}
	for name := range entries {
		if name == ".git/" || len(name) > 5 && name[:5] == ".git/" {
			t.Errorf(".git content leaked into context: %s", name)
		}
	}

	if hdr := entries[".nixpacks/assets/run.sh"]; hdr != nil && hdr.Mode&0o111 == 0 {
		t.Errorf("exec bit lost on run.sh: mode %o", hdr.Mode)
	}
	if hdr := entries["package.json"]; hdr != nil && (hdr.Uid != 0 || hdr.Gid != 0) {
		t.Errorf("ownership not zeroed: uid=%d gid=%d", hdr.Uid, hdr.Gid)
	}

	for _, link := range []string{"link.txt", "escape.txt"} {
		hdr := entries[link]
		if hdr == nil {
			t.Errorf("%s missing", link)
			continue
		}
		if hdr.Typeflag != tar.TypeSymlink {
			t.Errorf("%s typeflag = %v, want symlink (never followed)", link, hdr.Typeflag)
		}
		if bodies[link] != "" {
			t.Errorf("%s has file content — the link was followed", link)
		}
	}
	if entries["escape.txt"] != nil && entries["escape.txt"].Linkname != "/etc/hostname" {
		t.Errorf("escape.txt linkname = %q, want /etc/hostname", entries["escape.txt"].Linkname)
	}
}
