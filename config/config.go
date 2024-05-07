package config

import (
	"flag"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

/** ReadFile returns Configuration struct from a file path */
func ReadFile(path string) (*Configuration, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}
	var cfg = new(Configuration)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}

/** LoadConfigFromPath reads from path for local development */
func LoadConfigFromPath(pathToConfig string) Configuration {
	cfgPath := flag.String("p", pathToConfig, "Path to config file")
	flag.Parse()
	cfg, err := ReadFile(*cfgPath)
	if err != nil {
		panic(err.Error())
	}
	return *cfg
}

/** Gets an environment variable by key */
func GetEnvironmentVariable(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
