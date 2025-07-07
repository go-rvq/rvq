package env

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
	"github.com/theplant/osenv"
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

var InitialUserPasswordFile = osenv.Get("LOGIN_INITIAL_USER_PASSWORD_FILE", "Initial user password file", filepath.Join(PKG_DIR, ".initial_user_password"))

const DOT_ENV = ".env"

func init() {
	var (
		envFiles []string
		envf     = os.Getenv("ENV_FILES")
		envMap   = map[string]string{}
		load     = func(name string) (err error) {
			var content []byte
			if content, err = os.ReadFile(name); err != nil {
				return
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s: %v", name, err)
				}
			}()

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

	if envf == "" {
		envFiles = append(envFiles, filepath.Join(PKG_DIR, DOT_ENV))
	} else {
		envFiles = strings.Split(os.Getenv("ENV_FILES"), ",")
	}

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
			panic(err)
		}
	}

	for key, value := range envMap {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}
