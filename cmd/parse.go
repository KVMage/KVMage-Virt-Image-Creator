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
			  if key, val, ok := strings.Cut(line, "="); ok {
					  env[strings.TrimSpace(key)] = strings.TrimSpace(val)
			  }
	  }
	  return env, scanner.Err()
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
	  return re.ReplaceAllStringFunc(content, func(match string) string {
			  key := match[2 : len(match)-1]
			  if val, ok := env[key]; ok {
					  return val
			  }
			  return match
	  })
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

