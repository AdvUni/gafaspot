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
	SecretEngines []SecretEngineConfig //`yaml:"secretEngines"`
	Description   string               //`yaml:"description"`
}

// SecretEngineConfig is a struct to load information about one Secret Engine from config file.
type SecretEngineConfig struct {
	Name       string //`yaml:"name"`
	EngineType string `mapstructure:"type"`
	Role       string //`yaml:"role"`
}

type Environment struct {
	Name        string
	NamePlain   string
	HasSSH      bool
	Description string
}

type Reservation struct {
	ID      int
	Status  string
	User    string
	EnvName string
	Start   time.Time
	End     time.Time
	Subject string
	Labels  string
}
