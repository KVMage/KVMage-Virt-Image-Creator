package cmd

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func parseEnvFile(path string) (map[string]string, error) {
	env := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
			return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") {
					continue
			}

			line = strings.TrimPrefix(line, "export ")

			key, val, ok := strings.Cut(line, "=")
			if !ok {
					continue
			}

			key = strings.TrimSpace(key)
			val = strings.TrimSpace(val)

			if len(val) >= 2 {
					first, last := val[0], val[len(val)-1]
					if first == '"' && last == '"' {
							val = val[1 : len(val)-1]
							val = strings.ReplaceAll(val, `\n`, "\n")
							val = strings.ReplaceAll(val, `\t`, "\t")
							val = strings.ReplaceAll(val, `\\`, "\\")
							val = strings.ReplaceAll(val, `\"`, "\"")
					} else if first == '\'' && last == '\'' {
							val = val[1 : len(val)-1]
					} else {
							if idx := strings.Index(val, " #"); idx != -1 {
									val = strings.TrimSpace(val[:idx])
							}
					}
			} else {
					if idx := strings.Index(val, " #"); idx != -1 {
							val = strings.TrimSpace(val[:idx])
					}
			}

			env[key] = val
	}

	for key, val := range env {
			env[key] = expandEnvVars(val, env)
	}

	return env, scanner.Err()
}

func expandEnvVars(val string, env map[string]string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	val = re.ReplaceAllStringFunc(val, func(match string) string {
			key := match[2 : len(match)-1]
			if v, ok := env[key]; ok {
					return v
			}
			return match
	})

	re2 := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	val = re2.ReplaceAllStringFunc(val, func(match string) string {
			key := match[1:]
			if v, ok := env[key]; ok {
					return v
			}
			return match
	})

	return val
}

func loadEnvVars(configPath, envFilePath string) map[string]string {
	env := make(map[string]string)

	configDir := filepath.Dir(configPath)
	dotEnvPath := filepath.Join(configDir, ".env")
	if vars, err := parseEnvFile(dotEnvPath); err == nil {
			for k, v := range vars {
					env[k] = v
			}
	}

	if envFilePath != "" {
			if vars, err := parseEnvFile(envFilePath); err == nil {
					for k, v := range vars {
							env[k] = v
					}
			}
	}

	for _, e := range os.Environ() {
			if key, val, ok := strings.Cut(e, "="); ok {
					env[key] = val
			}
	}

	return env
}

func expandVars(content string, env map[string]string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		key := match[2 : len(match)-1]
		if val, ok := env[key]; ok {
			return val
		}
		return match
	})
	
	re2 := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	content = re2.ReplaceAllStringFunc(content, func(match string) string {
		key := match[1:]
		if val, ok := env[key]; ok {
			return val
		}
		return match
	})
	
	return content
}

func LoadConfig(path string) (map[string]*Options, error) {
	data, err := os.ReadFile(path)
	if err != nil {
			return nil, err
	}

	env := loadEnvVars(path, envFilePath)
	expanded := expandVars(string(data), env)

	var raw struct {
			KVMage map[string]*Options `yaml:"kvmage"`
	}

	if err := yaml.Unmarshal([]byte(expanded), &raw); err != nil {
			return nil, err
	}

	return raw.KVMage, nil
}
