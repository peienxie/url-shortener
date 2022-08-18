package config_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/peienxie/url-shortener/config"
	"github.com/stretchr/testify/require"
)

const envFilename = "app.env"

// create default app.env in current folder
var defaultEnv = map[string]string{
	"REDIS_ADDR":      "default_redis_addr",
	"API_SERVER_ADDR": "default_api_server_addr",
}

func createEnvFile(t *testing.T) {
	var b bytes.Buffer
	for k, v := range defaultEnv {
		b.WriteString(k + "=" + v + "\n")
	}

	err := os.WriteFile(envFilename, b.Bytes(), 0644)
	require.NoError(t, err)
}

func deleteEnvFile(t *testing.T) {
	err := os.Remove(envFilename)
	require.NoError(t, err)
}

func TestLoadConfigFromDotEnvFile(t *testing.T) {
	createEnvFile(t)

	config, err := config.LoadConfig(".")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	require.Equal(t, defaultEnv["REDIS_ADDR"], config.RedisAddr)
	require.Equal(t, defaultEnv["API_SERVER_ADDR"], config.ApiServerAddr)

	deleteEnvFile(t)
}

func TestOverrideConfigByEnvironmentVariables(t *testing.T) {
	os.Setenv("REDIS_ADDR", "env_var_redis_addr")
	os.Setenv("API_SERVER_ADDR", "env_var_api_server_addr")

	config, err := config.LoadConfig(".")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	require.Equal(t, "env_var_redis_addr", config.RedisAddr)
	require.Equal(t, "env_var_api_server_addr", config.ApiServerAddr)
}
