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

	contentBeginReservation = "Your reservation in Gafaspot for environment '%s' became active. You can login to retrieve your credentials."
	contentEndReservation   = "Your reservation in Gafaspot for environment '%s' expired. The credentials you received are not longer valid."

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
	logger        logging.Logger
	mailserver    string
	senderAddress string
)

func initMailing(l logging.Logger, config util.GafaspotConfig) {
	logger = l

	mailserver = "mail.au.de"
	senderAddress = "gafaspot@au.de"
}

func sendMail(recipient string, subject string, content string) error {

	msg := []byte(fmt.Sprintf(msgTemplate, senderAddress, recipient, subject, content))
	logger.Debugf("Assembled following email: %s", msg)

	err := smtp.SendMail(mailserver, nil, senderAddress, []string{recipient}, msg)

	return err
}

func sendBeginReservationMail(recipient string, creds util.ReservationCreds) {
	content := fmt.Sprintf(contentBeginReservation, creds.Env.NiceName)

	err := sendMail(recipient, subjectBeginReservation, content)
	if err != nil {
		logger.Errorf("failed to send mail to user %s at begin of reservation of env %s: %v", creds.Res.User, creds.Env.NiceName, err)
	}
}

func sendEndReservationMail(recipient string, reservation util.Reservation, envNiceName string) {
	content := fmt.Sprintf(contentEndReservation, envNiceName)

	err := sendMail(recipient, subjectEndReservation, content)
	if err != nil {
		logger.Errorf("failed to send mail to user %s at end of reservation of env %s: %v", reservation.User, envNiceName, err)
	}
}
