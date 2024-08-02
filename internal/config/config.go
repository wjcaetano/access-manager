package config

import (
	"fmt"
	"os"

	"github.com/magiconair/properties"
)

const (
	localConfigScope   = "resources/config/local.properties"
	localConfigFileEnv = "LOCAL_CONFIG_FILE_NAME"
	scopeEnv           = "SCOPE"
	appPathEnv         = "APP_PATH"
	localScope         = "local"
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
		AppPath  string `properties:"app_path,default="`
		Scope    string `properties:"scope,default="`
		Database Database
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
	configFile := getEnv(localConfigFileEnv, os.Getenv(appPathEnv)+localConfigScope)

	return properties.MustLoadFile(configFile, properties.UTF8)
}

func loadServiceProperties() (*properties.Properties, error) {
	return properties.LoadFile(os.Getenv(appPathEnv)+"resources/config/service.properties", properties.UTF8)
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
		return dns, fmt.Errorf("missing environment variables")
	}

	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", user, password, host, port, dbName)

	return dns, nil
}
