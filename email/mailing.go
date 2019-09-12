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

package email

import (
	"fmt"
	"net/smtp"

	"github.com/AdvUni/gafaspot/util"
	logging "github.com/alexcesaro/log"
)

const (

	// subjectBeginReservation and subjectEndReservation are the subjects Gafaspot
	// uses when mailing to its users.
	subjectBeginReservation = "Gafaspot notification: Reservation is active"
	subjectEndReservation   = "Gafaspot notification: Reservation expired"

	contentBeginReservation = `Your reservation in Gafaspot became active. You can login to retrieve your credentials.
Username: %s
Environment: %s
Reservation period: %s - %s
Reservation subject: %s

Environment Description ...
Creds ...`
	contentEndReservation = `Your reservation in Gafaspot expired. The credentials you received are not longer valid.
Username: %s
Environment: %s
Reservation period: %s - %s
Reservation subject: %s

`
	// msgTemplate is for creating RFC 822-style emails.
	// Following strings must be passed in the correct order:
	//   * the sender's mail address
	//   * the recipient's mail address
	//   * the subject
	//   * the content
	// The resulting message contains the e-mail headers From, To, and Subject
	msgTemplate = "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n"
)

var (
	// MailingEnabled is only true, if a mailserver is given in config. If not,
	// Gafaspot is not able to send mails anyway.
	MailingEnabled bool

	logger        logging.Logger
	mailserver    string
	senderAddress string
)

// InitMailing reads the email paramters from config and stores them as package variables.
func InitMailing(l logging.Logger, config util.GafaspotConfig) {
	logger = l

	mailserver = config.Mailserver
	if mailserver != "" {
		MailingEnabled = true
	}
	logger.Debugf("Mail server is specified: %v", MailingEnabled)

	senderAddress = config.GafaspotMailAddress
}

func sendMail(recipient string, subject string, content string) error {

	msg := []byte(fmt.Sprintf(msgTemplate, senderAddress, recipient, subject, content))
	logger.Debugf("Assembled following email: %s", msg)

	err := smtp.SendMail(mailserver, nil, senderAddress, []string{recipient}, msg)

	return err
}

// SendBeginReservationMail sends an e-mail to inform a user about the beginning of his reservation.
// recipient has to be the user's e-mail address.
func SendBeginReservationMail(recipient string, info util.ReservationCreds) {
	// TODO: improve content
	content := fmt.Sprintf(contentBeginReservation, info.Res.User, info.Env.NiceName, info.Res.Start.Format(util.TimeLayout), info.Res.End.Format(util.TimeLayout), info.Res.Subject)

	err := sendMail(recipient, subjectBeginReservation, content)
	if err != nil {
		logger.Errorf("failed to send mail to user %s at begin of reservation of env %s: %v", info.Res.User, info.Env.PlainName, err)
	}
}

// SendEndReservationMail sends an e-mail to inform a user about the end of his reservation.
// recipient has to be the user's e-mail address.
// As at a reservation's end there is no point in showing credentials, the info.Creds
// attribute is ignored and can be nil.
func SendEndReservationMail(recipient string, info util.ReservationCreds) {
	// TODO: improve content
	content := fmt.Sprintf(contentEndReservation, info.Res.User, info.Env.NiceName, info.Res.Start.Format(util.TimeLayout), info.Res.End.Format(util.TimeLayout), info.Res.Subject)

	err := sendMail(recipient, subjectEndReservation, content)
	if err != nil {
		logger.Errorf("failed to send mail to user %s at end of reservation of env %s: %v", info.Res.User, info.Env.PlainName, err)
	}
}
