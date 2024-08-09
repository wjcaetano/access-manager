package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
)

const (
	localConfigScope       = "resources/config/local.properties"
	applicationConfigScope = "resources/config/application.properties"
	LocalConfigFileEnv     = "LOCAL_CONFIG_FILE_NAME"
	scopeEnv               = "SCOPE"
	appPathEnv             = "APP_PATH"
	localScope             = "local"
)

type (
	Database struct {
		Cluster  string `properties:"cluster"`
		Name     string `properties:"name"`
		Password string `properties:"password"`
		Username string `properties:"username"`
		DSN      string `properties:"dsn"`
	}

	Configuration struct {
		AppPath  string   `properties:"app_path,default="`
		Scope    string   `properties:"scope,default="`
		Database Database `properties:"database"`
	}
)

func NewConfig() (Configuration, error) {
	prop, err := loadProperties()
	if err != nil {
		return Configuration{}, err
	}

	conf, err := decodeConfig(prop)
	if err != nil {
		return Configuration{}, err
	}

	conf.overrideConfigurations()

	dsn, err := buildDatabaseDSN()
	if err != nil {
		return Configuration{}, err
	}

	conf.Database.DSN = dsn

	return conf, nil
}

func loadProperties() (*properties.Properties, error) {
	if err := checkMandatoryEnvs(); err != nil {
		return nil, err
	}

	if getEnv(scopeEnv, "") == localScope {
		return loadLocalProperties(), nil
	}

	return loadServiceProperties()
}

func checkMandatoryEnvs() error {
	mandatoryEnvs := [...]string{appPathEnv, scopeEnv}
	app := getEnv(appPathEnv, "")
	fmt.Println(app)
	for _, env := range mandatoryEnvs {
		if _, found := os.LookupEnv(env); !found {
			return fmt.Errorf("environment %s not provided", env)
		}
	}

	return nil
}

// Configurations overwritten by environment variables.
func (c *Configuration) overrideConfigurations() {
	c.AppPath = getEnv(appPathEnv, "")
	c.Scope = getEnv(scopeEnv, "")
}

func getEnv(key string, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}

	return defaultValue
}

func loadLocalProperties() *properties.Properties {
	appPath, err := getProjectPath()
	if err != nil {
		return nil
	}
	configFile := filepath.Join(appPath, localConfigScope)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil
	}

	return properties.MustLoadFile(configFile, properties.UTF8)
}

func loadServiceProperties() (*properties.Properties, error) {
	inputConfig := os.Getenv("configFileName")
	if inputConfig == "" {
		inputConfig = applicationConfigScope
	}

	appPath, err := getProjectPath()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}

	configFile := filepath.Join(appPath, inputConfig)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file %s not found", inputConfig)
	}

	prop, _ := properties.LoadFile(configFile, properties.UTF8)

	return prop, nil
}

func decodeConfig(prop *properties.Properties) (Configuration, error) {
	var cfg Configuration
	if err := prop.Decode(&cfg); err != nil {
		return Configuration{}, err
	}

	return cfg, nil
}

func IsLocalScope() bool {
	return getEnv(scopeEnv, "") == localScope
}

func buildDatabaseDSN() (string, error) {
	var dns string

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
		return dns, fmt.Errorf("missing dns environment variables")
	}

	dns = strings.ToLower(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", user, password, host, port, dbName))

	return dns, nil
}

func getProjectPath() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting working directory: %w", err)
	}

	workingDir = filepath.Dir(filepath.Dir(workingDir))

	return workingDir, nil
}
