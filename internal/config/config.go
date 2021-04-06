package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mohitsinghs/wormholes/internal/database"
	"github.com/mohitsinghs/wormholes/internal/database/postgres"
	"github.com/mohitsinghs/wormholes/internal/database/redis"
	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_PORT = 3000
	DEFAULT_CONF = "config.yaml"
	DIR_PERM     = 0775
)

type Config struct {
	Port    int
	Backend string
	postgres.Postgres
	redis.Redis
}

func defaultConfig() Config {
	return Config{
		Port:     DEFAULT_PORT,
		Backend:  database.POSTGRES,
		Postgres: postgres.Default(),
		Redis:    redis.Default(),
	}
}

func configDir() string {
	home, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Error getting home dir : %v", err)
	}
	cfgDir := filepath.Join(home, "wh")
	_, err = os.Stat(cfgDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, DIR_PERM); err != nil {
			log.Printf("Error creating config dir : %v", err)
		}
	}
	return cfgDir
}

func writeDefault(cfgFile string) error {
	if f, err := os.Create(cfgFile); err != nil {
		return err
	} else {
		data, err := yaml.Marshal(defaultConfig())
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

func LoadDefault() (*Config, error) {
	cfgFile := filepath.Join(configDir(), "config.yaml")
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

func Merge(prefix string, conf *Config) {
	val := reflect.ValueOf(conf)
	set := map[string]interface{}{}
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		env := fmt.Sprintf("%s_%s", prefix, strings.ToUpper(f.Name))
		env = strings.Replace(env, "-", "_", -1)
		if v := os.Getenv(env); v != "" {
			if _, defined := set[f.Name]; !defined {
				flag.CommandLine.Set(f.Name, v)
			}
		} else {
			confv := reflect.Indirect(val).FieldByName(strings.Title(f.Name))
			if confv.IsValid() {
				flag.CommandLine.Set(f.Name, fmt.Sprint(confv))
			}
		}
	})
}
