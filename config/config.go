package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mohitsinghs/wormholes/constants"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     int
	Postgres Postgres
	Factory  FactoryConfig
	Pipe     PipeConfig
}

func Default() Config {
	return Config{
		Port:     constants.DEFAULT_PORT,
		Postgres: DefaultPostgres(),
		Factory:  DefaultFactory(),
		Pipe:     DefaultPipe(),
	}
}

func configDir() string {
	home, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Error getting home dir : %v", err)
	}
	cfgDir := filepath.Join(home, constants.DEFAULT_DIR)
	_, err = os.Stat(cfgDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, constants.DIR_PERM); err != nil {
			log.Printf("Error creating config dir : %v", err)
		}
	}
	return cfgDir
}

func writeDefault(cfgFile string) error {
	if f, err := os.Create(cfgFile); err != nil {
		return err
	} else {
		data, err := yaml.Marshal(Default())
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Load(cfgFile string) (*Config, error) {
	var conf *Config
	var err error
	if cfgFile != "" {
		conf, err = LoadFromFile(cfgFile)
	} else {
		conf, err = LoadDefault()
	}
	return conf, err
}

func LoadDefault() (*Config, error) {
	cfgFile := filepath.Join(configDir(), constants.DEFAULT_CONF)
	return LoadFromFile(cfgFile)
}

func LoadFromFile(cfgFile string) (*Config, error) {
	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			writeDefault(cfgFile)
		}
	}
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Printf("Error reading config file : %v", err)
		return nil, err
	}
	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		log.Printf("Error parsing config file : %v", err)
		return nil, err
	}
	return &conf, nil
}
