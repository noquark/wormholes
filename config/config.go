package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/mohitsinghs/wormholes/constants"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      int
	Postgres  Postgres
	Timescale Postgres
	Factory   FactoryConfig
	Pipe      PipeConfig
	Admin     Admin
}

func Load(cfgFile string) (*Config, error) {
	var conf *Config
	// Use default config path if no other path is passed
	if cfgFile == "" {
		cfgFile = filepath.Join(defaultConfigDir(), constants.DEFAULT_CONF)
	}
	info, err := os.Stat(cfgFile)
	// write and return config missing
	if os.IsNotExist(err) || info.Size() == 0 {
		conf = &Config{
			Port:      constants.DEFAULT_PORT,
			Postgres:  DefaultPostgres(),
			Timescale: DefaultTimescale(),
			Factory:   DefaultFactory(),
			Pipe:      DefaultPipe(),
			Admin:     DefaultAdmin(),
		}
		conf.Update(cfgFile)
		return conf, nil
	}
	// read and return config if exists
	if info.Mode().IsRegular() {
		data, err := os.ReadFile(cfgFile)
		if err != nil {
			log.Println("error reading config file : ", err)
			return nil, err
		}
		var conf Config
		if err := yaml.Unmarshal(data, &conf); err != nil {
			log.Println("Error parsing config file : ", err)
			return nil, err
		}
		return &conf, nil
	}
	if err != nil {
		log.Println("failed to read/write config : ", err)
		return nil, err
	} else {
		return nil, errors.New("failed to read/write config")
	}
}

// Update config on filesystem
func (c *Config) Update(cfgFile string) {
	log.Printf("Update called with config %v", c)
	data, err := yaml.Marshal(c)
	if err != nil {
		log.Println("config is not writable :", err)
		return
	}
	f, err := os.Create(cfgFile)
	if err != nil {
		log.Println("config is not writable :", err)
		return
	}
	_, err = f.Write(data)
	f.Close()
	if err != nil {
		log.Println("failed to update config : ", err)
	} else {
		log.Println("config updated")
	}
}

// Ensure and return default config directory
func defaultConfigDir() string {
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
