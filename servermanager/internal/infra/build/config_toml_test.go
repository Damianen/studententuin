package build

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"servermanager/internal/domain"

	toml "github.com/pelletier/go-toml/v2"
)

func readConfig(t *testing.T, dir string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		t.Fatalf("reading %s: %v", configFileName, err)
	}
	cfg := map[string]any{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parsing %s: %v", configFileName, err)
	}
	return cfg
}

func TestWriteBuildConfigNoOverrides(t *testing.T) {
	name, args, err := writeBuildConfig(domain.BuildInput{Dir: t.TempDir()})
	if err != nil {
		t.Fatalf("writeBuildConfig: %v", err)
	}
	if name != "" || args != nil {
		t.Errorf("got (%q, %v), want no config for an override-free input", name, args)
	}
}

func TestWriteBuildConfigOverrides(t *testing.T) {
	dir := t.TempDir()
	name, args, err := writeBuildConfig(domain.BuildInput{
		Dir:          dir,
		Env:          map[string]string{"FOO": "bar"},
		BuildCommand: "npm run build",
		StartCommand: "node server.js",
		Port:         3000,
	})
	if err != nil {
		t.Fatalf("writeBuildConfig: %v", err)
	}
	if name != configFileName {
		t.Errorf("config name = %q, want %q", name, configFileName)
	}

	cfg := readConfig(t, dir)
	vars := cfg["variables"].(map[string]any)
	if vars["FOO"] != "bar" || vars["PORT"] != "3000" {
		t.Errorf("variables = %v, want FOO=bar and injected PORT=3000", vars)
	}
	if cmd := cfg["start"].(map[string]any)["cmd"]; cmd != "node server.js" {
		t.Errorf("start.cmd = %v, want node server.js", cmd)
	}
	cmds := cfg["phases"].(map[string]any)["build"].(map[string]any)["cmds"].([]any)
	if len(cmds) != 1 || cmds[0] != "npm run build" {
		t.Errorf("phases.build.cmds = %v, want [npm run build]", cmds)
	}

	// Nixpacks turns variables into ARG+ENV and needs the values as build args.
	if len(args) != 2 || args["FOO"] == nil || *args["FOO"] != "bar" || args["PORT"] == nil || *args["PORT"] != "3000" {
		t.Errorf("build args = %v, want FOO=bar PORT=3000", args)
	}
}

func TestWriteBuildConfigUserPortWins(t *testing.T) {
	dir := t.TempDir()
	_, args, err := writeBuildConfig(domain.BuildInput{
		Dir:  dir,
		Env:  map[string]string{"PORT": "8081"},
		Port: 3000,
	})
	if err != nil {
		t.Fatalf("writeBuildConfig: %v", err)
	}
	if vars := readConfig(t, dir)["variables"].(map[string]any); vars["PORT"] != "8081" {
		t.Errorf("variables.PORT = %v, want the user's 8081", vars["PORT"])
	}
	if *args["PORT"] != "8081" {
		t.Errorf("build arg PORT = %v, want 8081", *args["PORT"])
	}
}

func TestWriteBuildConfigMergesRepoConfig(t *testing.T) {
	dir := t.TempDir()
	repoToml := []byte("[variables]\nKEEP = \"yes\"\nFOO = \"repo\"\n\n[phases.setup]\nnixPkgs = [\"nodejs\"]\n")
	if err := os.WriteFile(filepath.Join(dir, "nixpacks.toml"), repoToml, 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := writeBuildConfig(domain.BuildInput{
		Dir: dir,
		Env: map[string]string{"FOO": "request"},
	})
	if err != nil {
		t.Fatalf("writeBuildConfig: %v", err)
	}

	cfg := readConfig(t, dir)
	vars := cfg["variables"].(map[string]any)
	if vars["KEEP"] != "yes" {
		t.Errorf("repo-only variable dropped: %v", vars)
	}
	if vars["FOO"] != "request" {
		t.Errorf("request override lost: %v", vars)
	}
	if _, ok := cfg["phases"].(map[string]any)["setup"]; !ok {
		t.Errorf("repo phases.setup dropped: %v", cfg["phases"])
	}
}

func TestWriteBuildConfigRejectsControlChars(t *testing.T) {
	for _, cmd := range []string{"evil\ninjected", "tab\there", "null\x00byte"} {
		_, _, err := writeBuildConfig(domain.BuildInput{Dir: t.TempDir(), StartCommand: cmd})
		if !errors.Is(err, domain.ErrInvalid) {
			t.Errorf("StartCommand %q: err = %v, want ErrInvalid", cmd, err)
		}
	}
}

func TestWriteBuildConfigBadRepoToml(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "nixpacks.toml"), []byte("not [valid toml"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := writeBuildConfig(domain.BuildInput{Dir: dir, Port: 3000}); err == nil {
		t.Error("expected error for invalid repo nixpacks.toml")
	}
}
