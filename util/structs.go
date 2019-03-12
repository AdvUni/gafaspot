package util

import "time"

// GafaspotConfig is a struct to load every information from config file.
type GafaspotConfig struct {
	VaultAddress      string                       `mapstructure:"vault-address"`
	WebserviceAddress string                       `mapstructure:"webservice-address"`
	ApproleID         string                       `mapstructure:"approle-roleID"`
	ApproleSecret     string                       `mapstructure:"approle-secretID"`
	UserPolicy        string                       `mapstructure:"ldap-group-policy"`
	Database          string                       `mapstructure:"db-path"`
	Environments      map[string]EnvironmentConfig //`yaml:"environments"`
}

// EnvironmentConfig is a struct to load information about one environment from config file.
type EnvironmentConfig struct {
	NiceName      string               `mapstructure:"show-name"`
	Description   string               //`yaml:"description"`
	SecretEngines []SecretEngineConfig `mapstructure:"secret-engines"`
}

// SecretEngineConfig is a struct to load information about one Secret Engine from config file.
type SecretEngineConfig struct {
	NiceName   string `mapstructure:"name"`
	EngineType string `mapstructure:"type"`
	Role       string //`yaml:"role"`
}

type Environment struct {
	NiceName    string
	PlainName   string
	HasSSH      bool
	Description string
}

type Reservation struct {
	ID           int
	Status       string
	User         string
	EnvPlainName string
	Start        time.Time
	End          time.Time
	Subject      string
	Labels       string
}
