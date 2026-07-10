package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"servermanager/internal/domain"

	toml "github.com/pelletier/go-toml/v2"
)

// configFileName is the fixed filename the generator is pointed at — one
// constant literal, never derived from user input.
const configFileName = "stt-nixpacks.toml"

// writeBuildConfig materializes the request's env/build/start/port overrides
// as a nixpacks config file inside the checkout, merged over the repo's own
// nixpacks.toml so committed build config keeps working. Nixpacks declares
// each variable as ARG+ENV in the generated Dockerfile and expects the
// values as docker build args — those are returned for the daemon call.
// Returns ("", nil, nil) when the request carries no overrides.
func writeBuildConfig(in domain.BuildInput) (string, map[string]*string, error) {
	if err := domain.ValidateCommand("build command", in.BuildCommand); err != nil {
		return "", nil, err
	}
	if err := domain.ValidateCommand("start command", in.StartCommand); err != nil {
		return "", nil, err
	}

	variables := map[string]string{}
	for k, v := range in.Env {
		variables[k] = v
	}
	// The PaaS convention nixpacks apps listen on; an explicit user PORT wins.
	if in.Port > 0 {
		if _, ok := variables["PORT"]; !ok {
			variables["PORT"] = strconv.Itoa(in.Port)
		}
	}

	if len(variables) == 0 && in.BuildCommand == "" && in.StartCommand == "" {
		return "", nil, nil
	}

	cfg, err := loadRepoConfig(in.Dir)
	if err != nil {
		return "", nil, err
	}

	if len(variables) > 0 {
		vars, _ := cfg["variables"].(map[string]any)
		if vars == nil {
			vars = map[string]any{}
		}
		for k, v := range variables {
			vars[k] = v
		}
		cfg["variables"] = vars
	}
	if in.BuildCommand != "" {
		phases, _ := cfg["phases"].(map[string]any)
		if phases == nil {
			phases = map[string]any{}
		}
		buildPhase, _ := phases["build"].(map[string]any)
		if buildPhase == nil {
			buildPhase = map[string]any{}
		}
		buildPhase["cmds"] = []string{in.BuildCommand}
		phases["build"] = buildPhase
		cfg["phases"] = phases
	}
	if in.StartCommand != "" {
		start, _ := cfg["start"].(map[string]any)
		if start == nil {
			start = map[string]any{}
		}
		start["cmd"] = in.StartCommand
		cfg["start"] = start
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return "", nil, fmt.Errorf("encoding nixpacks config: %w", err)
	}
	// 0600: the config can carry env values and only the manager (and the
	// build context tar) needs to read it.
	if err := os.WriteFile(filepath.Join(in.Dir, configFileName), data, 0o600); err != nil {
		return "", nil, fmt.Errorf("writing nixpacks config: %w", err)
	}

	buildArgs := make(map[string]*string, len(variables))
	for k, v := range variables {
		buildArgs[k] = &v
	}
	return configFileName, buildArgs, nil
}

// loadRepoConfig reads the repository's own nixpacks.toml when present, so
// the request's overrides merge into it instead of silently replacing it.
func loadRepoConfig(dir string) (map[string]any, error) {
	cfg := map[string]any{}
	// #nosec G304 -- manager-generated temp dir + fixed filename.
	data, err := os.ReadFile(filepath.Join(dir, "nixpacks.toml"))
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading repo nixpacks.toml: %w", err)
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("repo nixpacks.toml is not valid TOML: %w", err)
	}
	return cfg, nil
}
