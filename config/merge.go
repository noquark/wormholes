package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Merge command line flags and environment variables with config.
func Merge(prefix string, conf *Config) {
	val := reflect.ValueOf(conf)
	set := map[string]interface{}{}

	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		env := fmt.Sprintf("%s_%s", prefix, strings.ToUpper(f.Name))
		env = strings.ReplaceAll(env, "-", "_")
		if v := os.Getenv(env); v != "" {
			if _, defined := set[f.Name]; !defined {
				_ = flag.CommandLine.Set(f.Name, v)
			}
		} else {
			confv := reflect.Indirect(val).FieldByName(strings.Title(f.Name))
			if confv.IsValid() {
				_ = flag.CommandLine.Set(f.Name, fmt.Sprint(confv))
			}
		}
	})
}
