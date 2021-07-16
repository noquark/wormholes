package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

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
