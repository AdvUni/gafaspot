// Copyright 2019, Advanced UniByte GmbH.
// Author Marie Lohbeck.
//
// This file is part of Gafaspot.
//
// Gafaspot is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gafaspot is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gafaspot.  If not, see <https://www.gnu.org/licenses/>.

package util

import (
	"html/template"
	"time"
)

// GafaspotConfig is a struct to load every information from config file.
type GafaspotConfig struct {
	WebserviceAddress string                       `mapstructure:"webservice-address"`
	MaxBookingDays    int                          `mapstructure:"max-reservation-duration-days"`
	MaxQueuingMonths  int                          `mapstructure:"max-queuing-time-months"`
	Database          string                       `mapstructure:"db-path"`
	DBTTLmonths       int                          `mapstructure:"database-ttl-months"`
	VaultAddress      string                       `mapstructure:"vault-address"`
	ApproleID         string                       `mapstructure:"approle-roleID"`
	ApproleSecret     string                       `mapstructure:"approle-secretID"`
	UserPolicy        string                       `mapstructure:"ldap-group-policy"`
	Environments      map[string]EnvironmentConfig //`yaml:"environments"`
}

// EnvironmentConfig is a struct to load information about one environment from config file.
type EnvironmentConfig struct {
	NiceName       string                `mapstructure:"show-name"`
	Description    string                //`yaml:"description"`
	SecretsEngines []SecretsEngineConfig `mapstructure:"secrets-engines"`
}

// SecretsEngineConfig is a struct to load information about one Secret Engine from config file.
type SecretsEngineConfig struct {
	NiceName   string `mapstructure:"name"`
	EngineType string `mapstructure:"type"`
	Role       string //`yaml:"role"`
}

type Environment struct {
	NiceName    string
	PlainName   string
	HasSSH      bool
	Description template.HTML
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
