package domain

import (
	"errors"
	"testing"
)

var testHosts = []string{"github.com", "gitlab.com"}

func TestValidateRepoURL(t *testing.T) {
	valid := []string{
		"https://github.com/user/repo",
		"https://github.com/user/repo.git",
		"https://GITHUB.COM/user/repo",
		" https://gitlab.com/group/sub/repo ",
	}
	for _, in := range valid {
		if err := ValidateRepoURL(in, testHosts); err != nil {
			t.Errorf("ValidateRepoURL(%q) = %v, want nil", in, err)
		}
	}

	invalid := []string{
		"",
		"http://github.com/user/repo",          // not https
		"git@github.com:user/repo.git",         // ssh
		"https://example.com/user/repo",        // host not allowed
		"https://www.github.com/user/repo",     // exact host match only
		"https://user:pass@github.com/u/r",     // credentials
		"https://github.com:8443/user/repo",    // explicit port
		"https://github.com/user/repo?x=1",     // query
		"https://github.com/user/repo#frag",    // fragment
		"https://github.com",                   // no repo path
		"https://github.com/",                  // no repo path
		"file:///etc/passwd",                   // local scheme
		"https://github.com.evil.io/user/repo", // host suffix trick
	}
	for _, in := range invalid {
		err := ValidateRepoURL(in, testHosts)
		if err == nil {
			t.Errorf("ValidateRepoURL(%q) = nil, want error", in)
			continue
		}
		if !errors.Is(err, ErrInvalid) {
			t.Errorf("ValidateRepoURL(%q) error %v does not wrap ErrInvalid", in, err)
		}
	}
}

func TestValidateBranch(t *testing.T) {
	valid := []string{"", "main", "feat/deploy-from-git", "release-1.2.3", "user_branch"}
	for _, in := range valid {
		if err := ValidateBranch(in); err != nil {
			t.Errorf("ValidateBranch(%q) = %v, want nil", in, err)
		}
	}

	invalid := []string{
		"has space", "semi;colon", "back`tick", "dollar$sign",
		"dot..dot", "-leading-dash", "/leading-slash", "branch.lock",
		string(make([]byte, 300)),
	}
	for _, in := range invalid {
		err := ValidateBranch(in)
		if err == nil {
			t.Errorf("ValidateBranch(%q) = nil, want error", in)
			continue
		}
		if !errors.Is(err, ErrInvalid) {
			t.Errorf("ValidateBranch(%q) error %v does not wrap ErrInvalid", in, err)
		}
	}
}

func TestAppImageRef(t *testing.T) {
	appID := "5f1c9450-9d6a-4b2e-8c8e-0a4f6f0c31b2"
	deployID := "a1b2c3d4-1111-2222-3333-444455556666"

	if got, want := AppImageRepo(appID), "stt-app-"+appID; got != want {
		t.Errorf("AppImageRepo = %q, want %q", got, want)
	}
	if got, want := AppImageRef(appID, deployID), "stt-app-"+appID+":a1b2c3d4"; got != want {
		t.Errorf("AppImageRef = %q, want %q", got, want)
	}
	if got, want := AppImageRef(appID, "abc"), "stt-app-"+appID+":abc"; got != want {
		t.Errorf("AppImageRef with short id = %q, want %q", got, want)
	}
}

func TestOnceCleanup(t *testing.T) {
	calls := 0
	fn := OnceCleanup(func() { calls++ })
	fn()
	fn()
	if calls != 1 {
		t.Errorf("cleanup ran %d times, want 1", calls)
	}
}
