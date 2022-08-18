package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration of application
type Config struct {
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	ApiServerAddr string `mapstructure:"API_SERVER_ADDR"`
}

// LoadConfig loads configuration from environment variables
// or search .env file from provided path
func LoadConfig(path string) (config Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	// the filename of configuration: app.env
	v.SetConfigName("app")
	v.SetConfigType("env")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = v.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
		default:
			return config, fmt.Errorf("read config err: %w", err)
		}
	}
	err = bindEnvs(v, config)
	if err != nil {
		return config, fmt.Errorf("bind environment variables err: %w", err)
	}

	// workaround because viper does not treat env vars the same as other config
	for _, key := range v.AllKeys() {
		val := v.Get(key)
		viper.Set(key, val)
	}

	err = v.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("unmarshal config err: %w", err)
	}

	return
}

// https://github.com/spf13/viper/issues/188
func bindEnvs(v *viper.Viper, iface interface{}, parts ...string) error {
	ifValue := reflect.ValueOf(iface)
	ifType := reflect.TypeOf(iface)
	for i := 0; i < ifType.NumField(); i++ {
		valueField := ifValue.Field(i)
		typeField := ifType.Field(i)
		tagValue, ok := typeField.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		var err error
		switch valueField.Kind() {
		case reflect.Struct:
			err = bindEnvs(v, valueField.Interface(), append(parts, tagValue)...)
		default:
			err = v.BindEnv(strings.Join(append(parts, tagValue), "."))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
