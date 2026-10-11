package env

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/go-rvq/rvq/x/features"
	"github.com/go-rvq/rvq/x/osenv"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var PKG_DIR = func() string {
	workDir := "."
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			var (
				cmd = exec.Command("go", "env", "GOMOD")
				out bytes.Buffer
			)
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				panic(err)
			}

			gomod := strings.TrimSpace(out.String())
			if gomod == "" {
				panic("GOMOD is empty")
			}
			return filepath.Dir(gomod)
		}
	}
	return workDir
}()

var InitialUserPasswordFile = osenv.Get("RVQ_LOGIN_INITIAL_USER_PASSWORD_FILE", "Initial user password file", filepath.Join(PKG_DIR, ".initial_user_password"))

// APP_YML is the configuration of the application: its environment (`env:`)
// and the features it turns on (`features:`, package features). APP_YAML is
// the same, by the other extension; with both, APP_YML is read.
const (
	APP_YML  = "app.yml"
	APP_YAML = "app.yaml"
)

// DOT_ENV is the configuration as it was: variables alone (godotenv). Read
// when there is no APP_YML.
const DOT_ENV = ".env"

// IsYAML says the file name is an app.yml (by its extension); else a .env.
func IsYAML(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".yml" || ext == ".yaml"
}

// ParseAppYAML is an app.yml: its `env:` — each value a string; a number or a
// bool as written; a list, its items joined by the separator of lists of
// paths (os.PathListSeparator, ":" on Linux), as PATH —, its `${VAR}`
// expanded from the environment; and its `features:`.
func ParseAppYAML(content []byte) (env map[string]string, feats map[string]any, err error) {
	var doc struct {
		Env      map[string]any `yaml:"env"`
		Features map[string]any `yaml:"features"`
	}
	if err = yaml.Unmarshal(content, &doc); err != nil {
		return nil, nil, fmt.Errorf("error parsing app.yml: %w", err)
	}
	env = make(map[string]string, len(doc.Env))
	for k, v := range doc.Env {
		var s string
		if s, err = envValue(v); err != nil {
			return nil, nil, fmt.Errorf("env.%s: %w", k, err)
		}
		env[k] = s
	}
	return env, doc.Features, nil
}

func envValue(v any) (string, error) {
	switch x := v.(type) {
	case nil:
		return "", nil
	case string:
		return envsubst.String(x)
	case []any:
		parts := make([]string, len(x))
		for i, item := range x {
			if _, nested := item.([]any); nested {
				return "", fmt.Errorf("a list in a list")
			}
			s, err := envValue(item)
			if err != nil {
				return "", err
			}
			parts[i] = s
		}
		return strings.Join(parts, string(os.PathListSeparator)), nil
	case map[string]any:
		return "", fmt.Errorf("a map is not a value")
	default:
		return fmt.Sprint(x), nil
	}
}

// files are the configuration files read: ENV_FILES ("," apart), else the
// app.yml (or app.yaml) in PKG_DIR, else its .env; required are those
// ENV_FILES names.
func files() (names []string, required bool) {
	if envf := os.Getenv("ENV_FILES"); envf != "" {
		return strings.Split(envf, ","), true
	}
	for _, name := range []string{APP_YML, APP_YAML, DOT_ENV} {
		p := filepath.Join(PKG_DIR, name)
		if _, err := os.Stat(p); err == nil {
			if name == DOT_ENV {
				fmt.Fprintf(os.Stderr, "WARNING: %s is deprecated: move its variables to the env: of %s\n", p, filepath.Join(PKG_DIR, APP_YML))
			}
			return []string{p}, false
		}
	}
	return nil, false
}

func init() {
	var (
		envMap = map[string]string{}
		load   = func(name string) (err error) {
			var content []byte
			if content, err = os.ReadFile(name); err != nil {
				return
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s: %v", name, err)
				}
			}()

			if IsYAML(name) {
				vars, feats, err := ParseAppYAML(content)
				if err != nil {
					return err
				}
				for k, v := range vars {
					envMap[k] = v
				}
				return features.Load(feats)
			}

			expanded, err := envsubst.String(string(content))

			if err != nil {
				return fmt.Errorf("error expanding dot env file: %w", err)
			}
			vars, err := godotenv.UnmarshalBytes([]byte(expanded))
			if err != nil {
				return fmt.Errorf("error parsing dot env file: %w", err)
			}
			for k, v := range vars {
				envMap[k] = v
			}
			return
		}
	)

	envFiles, required := files()

	if value := os.Getenv("LOGIN_INITIAL_USER_PASSWORD"); value == "" {
		if data, err := os.ReadFile(InitialUserPasswordFile); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		} else {
			data = bytes.TrimSpace(data)
			if len(data) > 0 {
				if err = os.Setenv("LOGIN_INITIAL_USER_PASSWORD", string(data)); err != nil {
					panic(err)
				}
			}
		}
	}

	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		kv := strings.SplitN(rawEnvLine, "=", 2)
		envMap[kv[0]] = kv[1]
	}

	for _, name := range envFiles {
		if name == "" {
			continue
		}

		if err := load(name); err != nil {
			if !required && os.IsNotExist(err) {
				continue
			}
			panic(err)
		}
	}

	for key, value := range envMap {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}
