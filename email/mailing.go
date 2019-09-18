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
	"bytes"
	"fmt"
	"net/smtp"
	"path"
	"text/template"
	"time"

	"github.com/AdvUni/gafaspot/util"
	logging "github.com/alexcesaro/log"
)

const (

	// subjectBeginReservation and subjectEndReservation are the subjects Gafaspot
	// uses when mailing to its users.
	subjectBeginReservation = "Gafaspot notification: Reservation is active"
	subjectEndReservation   = "Gafaspot notification: Reservation expired"

	// msgTemplate is for creating RFC 822-style emails.
	// Following strings must be passed in the correct order:
	//   * the sender's mail address
	//   * the recipient's mail address
	//   * the subject
	//   * the content
	// The resulting message contains the e-mail headers From, To, Subject and Content-Type (text/html)
	msgTemplate = "From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html\r\n\r\n%s\r\n"
)

var (
	// MailingEnabled is only true, if a mailserver is given in config. If not,
	// Gafaspot is not able to send mails anyway.
	MailingEnabled bool

	logger        logging.Logger
	mailserver    string
	senderAddress string

	startmailTmpl *template.Template
	endmailTmpl   *template.Template
)

// InitMailing reads the email paramters from config and stores them as package variables.
// Further it prepares the mail templates.
func InitMailing(l logging.Logger, config util.GafaspotConfig) {
	logger = l

	mailserver = config.Mailserver
	if mailserver != "" {
		MailingEnabled = true
	}
	logger.Debugf("Mail server is specified: %v", MailingEnabled)

	senderAddress = config.GafaspotMailAddress

	if MailingEnabled {
		const (
			startmailTmplFile = "email/templates/startmail.html"
			endmailTmplFile   = "email/templates/endmail.html"
		)
		var err error
		startmailTmpl, err = template.New(path.Base(startmailTmplFile)).Funcs(template.FuncMap{
			"formatDatetime": func(t time.Time) string { return t.Format(util.TimeLayout) },
		}).ParseFiles(startmailTmplFile)
		if err != nil {
			logger.Error(err)
		}
		endmailTmpl, err = template.New(path.Base(endmailTmplFile)).Funcs(template.FuncMap{
			"formatDatetime": func(t time.Time) string { return t.Format(util.TimeLayout) },
		}).ParseFiles(endmailTmplFile)
		if err != nil {
			logger.Error(err)
		}
	}
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
	var content bytes.Buffer
	err := startmailTmpl.Execute(&content, info)
	if err != nil {
		logger.Error(err)
	}
	err = sendMail(recipient, subjectBeginReservation, content.String())
	if err != nil {
		logger.Errorf("failed to send mail to user %s at begin of reservation of env %s: %v", info.Res.User, info.Env.PlainName, err)
	}
}

// SendEndReservationMail sends an e-mail to inform a user about the end of his reservation.
// recipient has to be the user's e-mail address.
// As at a reservation's end there is no point in showing credentials, the info.Creds
// attribute is ignored and can be nil.
func SendEndReservationMail(recipient string, info util.ReservationCreds) {
	var content bytes.Buffer
	err := endmailTmpl.Execute(&content, info)
	if err != nil {
		logger.Error(err)
	}
	err = sendMail(recipient, subjectEndReservation, content.String())
	if err != nil {
		logger.Errorf("failed to send mail to user %s at end of reservation of env %s: %v", info.Res.User, info.Env.PlainName, err)
	}
}
