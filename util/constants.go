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

const (
	// TimeLayout is the layout which must follow all time strings stored to and retrieved from database.
	// It can be used with time.ParseInLocation() and time.Format(). Such time strings are perfectly comparable.
	// As they don't contain any time zone information, it should be always parsed with time.Local.
	// So, each time value inside gafaspot are interpreted in the local time zone of the running server.
	TimeLayout = "2006-01-02 15:04"

	// SecEngTypes are constant strings to define the Secrets Engines Types supported by Gafaspot.

	// SecEngTypeAD is the type for Active Directory Secrets Engine.
	SecEngTypeAD = "ad"
	// SecEngTypeDB is the type for Database Secrets Engine.
	SecEngTypeDB = "database"
	// SecEngTypeOntap is the type for Ontap Secrets Engine.
	SecEngTypeOntap = "ontap"
	// SecEngTypeSSHPubkey is the type for SSH-Pubkey Secrets Engine.
	SecEngTypeSSHPubkey = "ssh-pubkey"
	// SecEngTypeSSH is the type for SSH Secrets Engine.
	SecEngTypeSSH = "ssh"
)
