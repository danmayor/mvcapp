package mvcapp

import (
	"bytes"
	"errors"
	"html/template"
	"net/mail"
	"os"
	"strings"

	gomail "gopkg.in/gomail.v2"
)

// EmailMessage is a simple wrapper providing common email object fields
// this object is used in the email connector for reading and writing
// email messages
type EmailMessage struct {
	// From email address
	From *mail.Address

	// To email addresses
	To []*mail.Address

	// CC email addresses
	CC []*mail.Address

	// BCC email addresses
	BCC []*mail.Address

	//Attachments is a slice of the file names to attach to this email
	Attachments []string

	// Headers are addition / optional email headers
	Headers map[string]string

	// Subject line of the email
	Subject string

	// Body or content of the email
	Body string
}

// EmailConnector is a simple wrapper around the gomail package allowing us to
// easily send and receive emails
type EmailConnector struct {
	Hostname string
	Port     int
	Username string
	Password string

	Sender gomail.Sender
}

// AddRecipient adds a new "To" recipient of this email message
func (emailMessage *EmailMessage) AddRecipient(to string) error {
	toAddr, err := mail.ParseAddress(to)
	if err != nil {
		return err
	}

	if emailMessage.To == nil {
		emailMessage.To = append([]*mail.Address{}, toAddr)
	} else {
		emailMessage.To = append(emailMessage.To, toAddr)
	}

	return nil
}

// AddCC adds a new "CC" or carbon copy recipient of this email message
func (emailMessage *EmailMessage) AddCC(cc string) error {
	ccAddr, err := mail.ParseAddress(cc)
	if err != nil {
		return err
	}

	if emailMessage.CC == nil {
		emailMessage.CC = append([]*mail.Address{}, ccAddr)
	} else {
		emailMessage.CC = append(emailMessage.CC, ccAddr)
	}

	return nil
}

// AddBCC adds a new "BCC" or blind carbon copy recipient of this email message
func (emailMessage *EmailMessage) AddBCC(bcc string) error {
	bccAddr, err := mail.ParseAddress(bcc)
	if err != nil {
		return err
	}

	if emailMessage.BCC == nil {
		emailMessage.BCC = append([]*mail.Address{}, bccAddr)
	} else {
		emailMessage.BCC = append(emailMessage.BCC, bccAddr)
	}

	return nil
}

// AddAttachment adds the provided filename to be attached to this email message
func (emailMessage *EmailMessage) AddAttachment(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return err
	}

	if emailMessage.Attachments == nil {
		emailMessage.Attachments = append([]string{}, filename)
	} else {
		emailMessage.Attachments = append(emailMessage.Attachments, filename)
	}

	return nil
}

// NewEmailMessage constructs a new EmailMessage object from the provided arguments
func NewEmailMessage(from string, to string, subject string, body string) (*EmailMessage, error) {
	fromAddress, err := mail.ParseAddress(from)
	if err != nil {
		return nil, err
	}

	toAddress, err := mail.ParseAddress(to)
	if err != nil {
		return nil, err
	}

	return &EmailMessage{
		From:    fromAddress,
		To:      append([]*mail.Address{}, toAddress),
		Subject: subject,
		Body:    body,
	}, nil
}

// NewEmailFromTemplate executes the provided templatePath and data model to constuct the body
// text. This and other provided values are then used to call NewEmailMessage
func NewEmailFromTemplate(from string, to string, subject string, templatePath string, model interface{}) (*EmailMessage, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"RawHTML": RawHTML,
	}

	t, err := template.New("EmailMessage").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	var content bytes.Buffer
	if err := t.Execute(&content, model); err != nil {
		return nil, err
	}

	return NewEmailMessage(from, to, subject, content.String())
}

// NewEmailConnector constructs a new EmailConnector object from the provided arguments
func NewEmailConnector(hostname string, port int, username string, password string) *EmailConnector {
	return &EmailConnector{
		Hostname: hostname,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// SendMail sends the provided email message through the connector object
func (connector *EmailConnector) SendMail(emailMessage *EmailMessage) error {
	if connector.Sender == nil {
		dialer := gomail.NewDialer(connector.Hostname, connector.Port, connector.Username, connector.Password)
		sender, err := dialer.Dial()

		if err != nil {
			return err
		}

		connector.Sender = sender
	}

	if emailMessage.To == nil || len(emailMessage.To) <= 0 {
		return errors.New("No recipients provided, can not send email")
	}

	message := gomail.NewMessage()
	message.SetAddressHeader("From", emailMessage.From.Address, emailMessage.From.Name)
	message.SetHeader("Subject", emailMessage.Subject)
	message.SetBody("text/html", emailMessage.Body)

	toString := []string{}
	for _, to := range emailMessage.To {
		toString = append(toString, to.String())
	}
	message.SetHeader("To", toString...)

	if emailMessage.CC != nil && len(emailMessage.CC) > 0 {
		toString = []string{}
		for _, cc := range emailMessage.CC {
			toString = append(toString, cc.String())
		}
		message.SetHeader("Cc", toString...)
	}

	if emailMessage.BCC != nil && len(emailMessage.BCC) > 0 {
		toString = []string{}
		for _, bcc := range emailMessage.BCC {
			toString = append(toString, bcc.String())
		}
		message.SetHeader("Bcc", toString...)
	}

	if emailMessage.Attachments != nil && len(emailMessage.Attachments) > 0 {
		for _, att := range emailMessage.Attachments {
			message.Attach(att)
		}
	}

	if err := gomail.Send(connector.Sender, message); err != nil {
		return err
	}

	return nil
}
