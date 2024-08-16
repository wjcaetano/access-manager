package config

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/magiconair/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DecodeConfig(t *testing.T) {
	prop := properties.LoadMap(map[string]string{
		"scope":             "local",
		"database.cluster":  "db",
		"database.name":     "testlocal",
		"database.password": "root",
		"database.username": "root",
		"database.dsn":      "user:password@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4&loc=Local",
	})
	t.Run("should return success when decode config", func(t *testing.T) {
		expectedResult := Configuration{
			Scope:    "local",
			Database: buildDatabase(),
		}

		result, err := decodeConfig(prop)
		require.NoError(t, err)

		isEqual := reflect.DeepEqual(result, expectedResult)
		assert.True(t, isEqual)
	})

	t.Run("should return error when decode config", func(t *testing.T) {
		prop := properties.LoadMap(map[string]string{
			"scope": "local",
		})
		_, err := decodeConfig(prop)
		require.Error(t, err)
	})
}

func Test_GetEnv(t *testing.T) {
	defaultEnv := setupTestGetEnv(t)
	defer defaultEnv(t)
	t.Run("should given an loaded env return the value", func(t *testing.T) {
		expectedResult := "any env value"
		result := getEnv("TEST_ENV_KEY", "")
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should given an unloaded env return the default value", func(t *testing.T) {
		expectedResult := "default value"
		result := getEnv("ANY_ENV_KEY", "default value")
		assert.Equal(t, expectedResult, result)
	})
}

func Test_OverrideConfigurations(t *testing.T) {
	defaultEnv := setupTestOverrideConfigurations(t)
	defer defaultEnv(t)

	t.Run("should given a valid configuration and env with other appPath value should override appPath property", func(t *testing.T) {
		expectedResult := Configuration{
			AppPath:  "PATH",
			Scope:    "local",
			Database: Database{},
		}
		config := Configuration{}
		config.overrideConfigurations()

		isEqual := reflect.DeepEqual(config, expectedResult)
		assert.True(t, isEqual)
	})
}

func Test_LoadProperties(t *testing.T) {
	t.Run("should return error when mandatory envs not provided", func(t *testing.T) {
		cleanEnvsNewConfig(t)
		expectedError := "environment APP_PATH not provided"
		_, err := loadProperties()
		require.EqualError(t, err, expectedError)
	})

	t.Run("should return local properties when scope is local", func(t *testing.T) {
		cleanEnvsNewConfig(t)
		setupTesLoadLocalProperties(t)

		result, err := loadProperties()
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should return service properties when scope is not local", func(t *testing.T) {
		cleanEnvsNewConfig(t)
		setupSubTestLoadServiceProperties(t)

		result, err := loadProperties()
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func Test_NewConfig(t *testing.T) {
	t.Run("given local empty scope and app_path should return error", func(t *testing.T) {
		cleanEnvsNewConfig(t)

		expectedError := fmt.Errorf("environment APP_PATH not provided")
		expectedResult := Configuration{}
		os.Setenv("LOCAL_CONFIG_FILE_NAME", "testdata/valid.properties")
		result, err := NewConfig()

		assert.Equal(t, expectedError, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("given a valid local properties should load properties", func(t *testing.T) {
		cleanEnvsNewConfig(t)
		setupNewConfigEnv(t)

		appPath, err := getProjectPath()
		require.NoError(t, err)

		expectedResult := Configuration{
			AppPath: appPath,
			Scope:   "",
			Database: Database{
				Cluster:  "db",
				Name:     "testlocal",
				Password: "root",
				Username: "root",
				DSN:      "root:root@tcp(testlocal:1234)/testlocal?parseTime=true&charset=utf8mb4&loc=Local",
			},
		}

		result, err := NewConfig()

		require.NoError(t, err)
		require.Equal(t, expectedResult, result)
	})
}

func setupTesLoadLocalProperties(t *testing.T) func(t *testing.T) {
	t.Log("setup test load local properties")

	appPath, err := getProjectPath()
	require.NoError(t, err)

	_ = os.Setenv("APP_PATH", appPath)
	_ = os.Setenv("SCOPE", localConfigScope)

	return func(t *testing.T) {
		t.Log("teardown test load local properties")
		_ = os.Unsetenv("SCOPE")
		_ = os.Unsetenv("APP_PATH")
	}

}

func setupSubTestLoadServiceProperties(t *testing.T) func(t *testing.T) {
	t.Log("setup sub test load service properties")
	appPath, err := getProjectPath()
	require.NoError(t, err)

	_ = os.Setenv("SCOPE", "")
	_ = os.Setenv("APP_PATH", appPath)

	return func(t *testing.T) {
		t.Log("teardown sub test load service properties")
		_ = os.Unsetenv("SCOPE")
		_ = os.Unsetenv("APP_PATH")
	}

}

func setupTestOverrideConfigurations(t *testing.T) func(t *testing.T) {
	t.Log("setup test override configurations")
	os.Setenv("APP_PATH", "PATH")
	os.Setenv("SCOPE", "local")

	return func(t *testing.T) {
		t.Log("teardown test override configurations")
		os.Unsetenv("APP_PATH")
		os.Unsetenv("SCOPE")
	}
}

func buildDatabase() Database {
	return Database{
		Cluster:  "db",
		Name:     "testlocal",
		Password: "root",
		Username: "root",
		DSN:      "user:password@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4&loc=Local",
	}
}

func cleanEnvsNewConfig(t *testing.T) {
	t.Log("clean envs new config")
	_ = os.Unsetenv("LOCAL_CONFIG_FILE_NAME")
	_ = os.Unsetenv("SCOPE")
	_ = os.Unsetenv("APP_PATH")
}

func setupTestGetEnv(t *testing.T) func(t *testing.T) {
	t.Log("setup test get env")
	_ = os.Setenv("TEST_ENV_KEY", "any env value")

	return func(t *testing.T) {
		t.Log("teardown test get env")
		_ = os.Unsetenv("TEST_ENV_KEY")
	}
}

func setupNewConfigEnv(t *testing.T) {
	t.Log("setup new config env")
	appPath, _ := getProjectPath()

	_ = os.Setenv("configFileName", applicationConfigScope)
	_ = os.Setenv("SCOPE", "")
	_ = os.Setenv("APP_PATH", appPath)
	_ = os.Setenv("DB_USER", "root")
	_ = os.Setenv("DB_PASSWORD", "root")
	_ = os.Setenv("DB_PORT", "1234")
	_ = os.Setenv("DB_NAME", "testlocal")
	_ = os.Setenv("DB_CLUSTER", "db")
	_ = os.Setenv("DB_HOST", "testlocal")

	dsn, err := buildDatabaseDSN()
	require.NoError(t, err)

	_ = os.Setenv("DB_DSN", dsn)
}
