package build

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// writeTar streams dir as the docker build context. Top-level .git is
// skipped (history is neither needed in the image nor cheap to send),
// symlinks are emitted as links and never followed — a followed link could
// smuggle host files into a user image — and ownership is zeroed so the
// context is deterministic. File opens go through an os.Root scoped to the
// clone dir, so no rename/symlink race can reach outside it.
func writeTar(w io.Writer, dir string) error {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("opening context root: %w", err)
	}
	defer func() { _ = root.Close() }()

	tw := tar.NewWriter(w)

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		name := filepath.ToSlash(rel)
		if name == ".git" || strings.HasPrefix(name, ".git/") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		var link string
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		case !info.Mode().IsRegular() && !info.IsDir():
			return nil // sockets, devices, pipes: not build context material
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return fmt.Errorf("tar header for %s: %w", name, err)
		}
		hdr.Name = name
		if info.IsDir() {
			hdr.Name += "/"
		}
		hdr.Uid, hdr.Gid = 0, 0
		hdr.Uname, hdr.Gname = "", ""

		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("writing tar header %s: %w", name, err)
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := root.Open(name)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, f)
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return fmt.Errorf("streaming %s into context: %w", name, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return tw.Close()
}
