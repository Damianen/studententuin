package domain

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// branchRe is deliberately far stricter than git's own ref rules: repo owners
// control branch names, but the value still travels into clone options.
var branchRe = regexp.MustCompile(`^[A-Za-z0-9._/-]{1,255}$`)

// ValidateRepoURL enforces the fetch policy from §3.4: https only, host on the
// server-side allowlist, and nothing beyond a plain repository path — no
// userinfo, custom port, query, or fragment.
func ValidateRepoURL(raw string, allowedHosts []string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return fmt.Errorf("repository url: %w", ErrInvalid)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("repository url must use https: %w", ErrInvalid)
	}
	if u.User != nil {
		return fmt.Errorf("repository url must not contain credentials: %w", ErrInvalid)
	}
	if u.Port() != "" {
		return fmt.Errorf("repository url must not set a port: %w", ErrInvalid)
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("repository url must not contain a query or fragment: %w", ErrInvalid)
	}
	if strings.TrimLeft(u.Path, "/") == "" {
		return fmt.Errorf("repository url is missing a repository path: %w", ErrInvalid)
	}
	host := strings.ToLower(u.Hostname())
	for _, allowed := range allowedHosts {
		if host == strings.ToLower(strings.TrimSpace(allowed)) {
			return nil
		}
	}
	return fmt.Errorf("repository host %q is not allowed (allowed: %s): %w",
		host, strings.Join(allowedHosts, ", "), ErrInvalid)
}

// ValidateBranch accepts an empty branch (= the remote's default branch) or a
// conservative branch-name charset without traversal-looking segments.
func ValidateBranch(branch string) error {
	if branch == "" {
		return nil
	}
	if !branchRe.MatchString(branch) {
		return fmt.Errorf("branch %q contains unsupported characters: %w", branch, ErrInvalid)
	}
	if strings.Contains(branch, "..") || strings.HasPrefix(branch, "-") || strings.HasPrefix(branch, "/") ||
		strings.HasSuffix(branch, ".lock") {
		return fmt.Errorf("branch %q is not a valid ref name: %w", branch, ErrInvalid)
	}
	return nil
}

// SourceCheckout is a fetched working tree in a throwaway directory. Cleanup
// removes the directory and is safe to call multiple times.
type SourceCheckout struct {
	Dir       string
	CommitSHA string
	Cleanup   func()
}

// OnceCleanup wraps fn so concurrent/duplicate Cleanup calls run it once.
func OnceCleanup(fn func()) func() {
	var once sync.Once
	return func() { once.Do(fn) }
}
