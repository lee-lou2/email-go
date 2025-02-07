package config_test

import (
	"aws-ses-sender-go/config"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var envOnce sync.Once

func TestGetEnv_ExistingVariable(t *testing.T) {
	key := "TEST_EXISTING_VARIABLE"
	expectedValue := "test_value"
	os.Setenv(key, expectedValue)

	defer os.Unsetenv(key) // Ensure cleanup after test

	value := config.GetEnv(key)
	assert.Equal(t, expectedValue, value, "The environment variable value should match the expected value")
}

func TestGetEnv_NonExistingVariable(t *testing.T) {
	key := "TEST_NON_EXISTING_VARIABLE"

	value := config.GetEnv(key)
	assert.Equal(t, "", value, "Non-existing environment variable should return an empty string")
}

func TestGetEnv_WithDefaultValue(t *testing.T) {
	key := "TEST_NON_EXISTING_VARIABLE"
	defaultValue := "default_value"

	value := config.GetEnv(key, defaultValue)
	assert.Equal(t, defaultValue, value, "Non-existing environment variable should return the default value")
}

func TestGetEnv_ExistingVariableWithDefaultValue(t *testing.T) {
	key := "TEST_EXISTING_VARIABLE_WITH_DEFAULT"
	expectedValue := "existing_value"
	defaultValue := "default_value"
	os.Setenv(key, expectedValue)

	defer os.Unsetenv(key) // Ensure cleanup after test

	value := config.GetEnv(key, defaultValue)
	assert.Equal(t, expectedValue, value, "Existing environment variable should override the default value")
}

func TestGetEnv_LoadEnvOnce(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_env")
	assert.NoError(t, err, "Temporary directory creation should not fail")
	defer os.RemoveAll(tempDir)

	originalDir, err := os.Getwd()
	assert.NoError(t, err, "Fetching the current working directory should not fail")

	err = os.Chdir(tempDir)
	assert.NoError(t, err, "Changing to the temporary directory should not fail")
	defer os.Chdir(originalDir)

	envOnce.Do(func() {})

	config.GetEnv("TEST_VARIABLE")

	// Manually test log output for `.env` file warning if required
}
