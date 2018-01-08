/*
	Digivance MVC Application Framework - Unit Tests
	Email Connector Feature Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.3.0 compatibility of emailconnector.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in emailconnector.go
*/

package mvcapp_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/digivance/mvcapp"
)

// Because these test functions all require credentials, you'll need to modify the tests and
// set this to true, sorry. But all have been tested :)
const (
	doEmailConnectorTests = false
	emailHostname         = "smtp.your-domain.com"
	emailUsername         = "you@your-domain.com"
	emailPassword         = "your-secret-password"

	emailTo  = "to@your-domain.com"
	emailCC  = "cc@your-domain.com"
	emailBCC = "bcc@your-domain.com"

	emailFrom = "you@your-domain.com"
)

// TestEmailMessage_AddRecipient ensures that the EmailMessage.AddRecipient method works as expected
func TestEmailMessage_AddRecipient(t *testing.T) {
	msg, err := mvcapp.NewEmailMessage(emailFrom, emailTo, "Subject", "Body")
	if err != nil {
		t.Fatalf("Failed to test adding recipients: %s", err)
	}

	if err = msg.AddRecipient("fail me"); err == nil {
		t.Errorf("Failed to test failing to add invalid recipient")
	}

	msg.To = nil
	if err = msg.AddRecipient(emailTo); err != nil {
		t.Errorf("Failed to test adding recipient when none exist: %s", err)
	}
}

// TestEmailMessage_AddCC ensures that the EmailMessage.AddCC method works as expected
func TestEmailMessage_AddCC(t *testing.T) {
	msg, err := mvcapp.NewEmailMessage(emailFrom, emailTo, "Subject", "Body")
	if err != nil {
		t.Fatalf("Failed to test adding carbon copy (CC) recipients: %s", err)
	}

	if err = msg.AddCC("fail me"); err == nil {
		t.Errorf("Failed to test failing to add invalid cc recipient")
	}

	if err = msg.AddCC(emailTo); err != nil {
		t.Errorf("Failed to test adding cc recipient when none exist: %s", err)
	}

	if err = msg.AddCC("fail me"); err == nil {
		t.Errorf("Failed to test failing to add invalid cc recipient")
	}

	if err = msg.AddCC(emailTo); err != nil {
		t.Errorf("Failed to test adding cc recipient when none exist: %s", err)
	}
}

// TestEmailMessage_AddBCC ensures that the EmailMessage.AddBCC method works as expected
func TestEmailMessage_AddBCC(t *testing.T) {
	msg, err := mvcapp.NewEmailMessage(emailFrom, emailTo, "Subject", "Body")
	if err != nil {
		t.Fatalf("Failed to test adding blind carbon copy (BCC) recipients: %s", err)
	}

	if err = msg.AddBCC("fail me"); err == nil {
		t.Errorf("Failed to test failing to add invalid bcc recipient")
	}

	if err = msg.AddBCC(emailTo); err != nil {
		t.Errorf("Failed to test adding bcc recipient when none exist: %s", err)
	}

	if err = msg.AddBCC("fail me"); err == nil {
		t.Errorf("Failed to test failing to add invalid bcc recipient")
	}

	if err = msg.AddBCC(emailTo); err != nil {
		t.Errorf("Failed to test adding bcc recipient when none exist: %s", err)
	}
}

// TestEmailMessage_AddAttachment ensures that the EmailMessage.AddAttachment works as expected
func TestEmailMessage_AddAttachment(t *testing.T) {
	filename := mvcapp.GetApplicationPath() + "/testfile.txt"
	msg, err := mvcapp.NewEmailMessage(emailFrom, emailTo, "Subject", "Body")
	if err != nil {
		t.Fatalf("Failed to test adding attachment to email message: %s", err)
	}

	if err = msg.AddAttachment(filename); err == nil {
		t.Error("Failed to fail when attaching file that doesn't exist")
	}

	if err = ioutil.WriteFile(filename, []byte("Hello World"), 0644); err != nil {
		t.Errorf("Failed to write file to test attachments: %s", err)
	}
	defer os.RemoveAll(filename)

	if err = msg.AddAttachment(filename); err != nil {
		t.Errorf("Failed to attach file: %s", err)
	}

	if err = msg.AddAttachment(filename); err != nil {
		t.Errorf("Failed to attach file: %s", err)
	}
}

// TestNewEmailMessage ensures that the mvcapp.NewEmailMessage returns the expected value
func TestNewEmailMessage(t *testing.T) {
	_, err := mvcapp.NewEmailMessage("failer", emailTo, "subject", "body")
	if err == nil {
		t.Errorf("Failed to test creating a new email message: %s", err)
	}
}

// TestNewEmailMessageFromTemplate ensures that the mvcapp.NewEmailMessageFromTemplate
// returns the expected value
func TestNewEmailMessageFromTemplate(t *testing.T) {
	templateData := []byte("{{ define \"EmailMessage\" }} Hello {{ .Failure }}! {{ end }}")
	filename := mvcapp.GetApplicationPath() + "/testmail.tpl"
	err := ioutil.WriteFile(filename, templateData, 0644)
	defer os.RemoveAll(filename)
	if err != nil {
		t.Fatalf("Failed to create new email template for testing: %s", err)
	}

	msg, err := mvcapp.NewEmailMessageFromTemplate(emailFrom, emailTo, "Subject", filename, "failme")
	if err == nil {
		t.Errorf("Failed to fail if template fails... lol: %s", msg.Body)
	}
}

// TestEmailConnector_SendMail ensures that the EmailConnector.SendMail method operates as expected
// Note you will need to modify this method with your email credentials / recipients and set doEmailConnectorTests
// to true to include this with your unit tests.
func TestEmailConnector_SendMail(t *testing.T) {
	if !doEmailConnectorTests {
		return
	}

	con := mvcapp.NewEmailConnector(emailHostname, 587, emailUsername, emailPassword)
	msg, err := mvcapp.NewEmailMessage(emailFrom, emailTo, "Testing Email", "<strong>Test email</strong><br /><br /> This email was generated via Golang unit tests to proof ensure that the MvcApp Email Connector is working as expected.")

	if err != nil {
		t.Error(err)
	}

	err = con.SendMail(msg)
	if err != nil {
		t.Error(err)
	}

	msg.To = nil
	err = con.SendMail(msg)
	if err == nil {
		t.Error("Failed to fail sending to missing recipients")
	}

	con = mvcapp.NewEmailConnector("fail.mail.com", 587, emailUsername, emailPassword)
	err = con.SendMail(msg)
	if err == nil {
		t.Error("Failed to fail sending through invalid host")
	}

	_, err = mvcapp.NewEmailMessage(emailFrom, "failme as invalid to address", "Subject", "Body")
	if err == nil {
		t.Error("Failed to prevent creation of email message object to invalid recipient!")
	}
}

// TestEmailConnector_SendTemplateMail ensures that the EmailConnector.SendTemplateMail method operates
// as expected Note you will need to modify this method with your email credentials / recipients and set
// doEmailConnectorTests to true to include this with your unit tests.
func TestEmailConnector_SendTemplateMail(t *testing.T) {
	if !doEmailConnectorTests {
		return
	}

	con := mvcapp.NewEmailConnector(emailHostname, 587, emailUsername, emailPassword)
	data := struct {
		CompanyName string
	}{
		CompanyName: "Digivance Technologies",
	}

	tdata := "{{ define \"EmailMessage\" }}<h1>Template Test</h1>The maker of this library is {{ .CompanyName }}!<br /><br />{{ end }}"
	tpath := fmt.Sprintf("%s/testMail.htm", mvcapp.GetApplicationPath())
	ioutil.WriteFile(tpath, []byte(tdata), 0644)
	defer os.RemoveAll(tpath)

	msg, err := mvcapp.NewEmailMessageFromTemplate(emailFrom, emailTo, "Template Test!", tpath, data)

	if err != nil {
		t.Error(err)
	}

	err = con.SendMail(msg)
	if err != nil {
		t.Error(err)
	}

}

// TestEmailConnector_SendWithAttachment ensures that the EmailConnector.SendWithAttachment method operates
// as expected Note you will need to modify this method with your email credentials / recipients and set
// doEmailConnectorTests to true to include this with your unit tests.
func TestEmailConnector_SendWithAttachment(t *testing.T) {
	if !doEmailConnectorTests {
		return
	}

	con := mvcapp.NewEmailConnector(emailHostname, 587, emailUsername, emailPassword)
	msg, err := mvcapp.NewEmailMessage(emailFrom, emailTo, "Testing Email", "<strong>Test email with CC and BCC recipients!</strong><br /><br /> This email was generated via Golang unit tests to proof ensure that the MvcApp Email Connector is working as expected.")

	if err != nil {
		t.Error(err)
	}

	filename := fmt.Sprintf("%s/attachment.txt", mvcapp.GetApplicationPath())
	defer os.RemoveAll(filename)
	if err = ioutil.WriteFile(filename, []byte("A file"), 0644); err != nil {
		t.Errorf("Failed to create temporary file for attachment: %s", err)
	}

	if err = msg.AddAttachment(filename); err != nil {
		t.Error(err)
	}

	err = con.SendMail(msg)
	if err != nil {
		t.Error(err)
	}
}

// TestEmailConnector_SendMultiRecipients ensures that the EmailConnector.SendMultiRecipients method operates
// as expected Note you will need to modify this method with your email credentials / recipients and set
// doEmailConnectorTests to true to include this with your unit tests.
func TestEmailConnector_SendMultiRecipients(t *testing.T) {
	if !doEmailConnectorTests {
		return
	}

	con := mvcapp.NewEmailConnector(emailHostname, 587, emailUsername, emailPassword)
	msg, err := mvcapp.NewEmailMessage(emailFrom, emailTo, "Testing Email", "<strong>Test email with CC and BCC recipients!</strong><br /><br /> This email was generated via Golang unit tests to proof ensure that the MvcApp Email Connector is working as expected.")

	if err != nil {
		t.Error(err)
	}

	if err = msg.AddRecipient(emailTo); err != nil {
		t.Error(err)
	}

	if err = msg.AddCC(emailCC); err != nil {
		t.Error(err)
	}

	if err = msg.AddBCC(emailBCC); err != nil {
		t.Error(err)
	}

	err = con.SendMail(msg)
	if err != nil {
		t.Error(err)
	}
}
