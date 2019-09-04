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

// Environment is a struct to store the information of one row from database table environments.
// The Description is of type template.HTML, as this type will not be escaped when served with a
// golang http.Template. This enables the gafaspot config writer to put some HTML code inside the
// descriptions for the environments.
type Environment struct {
	NiceName    string
	PlainName   string
	HasSSH      bool
	Description template.HTML
}

// Reservation is a struct to store the information of one row from database table reservations.
// (only database column delete_on is not included).
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

// ReservationCreds is a struct to bundle up credentials for a reservation. ReservationCreds
// can hold the credentials itself, the Environment, they belong to, and the associated Reservation,
// for which the credentials were created.
// The Creds attribute is a map to store one map for each Secrets Engine, which contains some
// key-value pairs as they are retrieved by a KV Secrets Engines.
type ReservationCreds struct {
	Res   Reservation
	Env   Environment
	Creds map[string]map[string]interface{}
}
