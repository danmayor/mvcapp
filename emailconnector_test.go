/*
	Digivance MVC Application Framework - Unit Tests
	Email Connector Feature Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.2.0 compatibility of emailconnector.go functions. These functions are written
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
var doTests = false

// TestEmailConnector_SendMail ensures that the EmailConnector.SendMail method operates as expected
// Note you will need to modify this method with your email credentials / recipients and set doTests
// to true to include this with your unit tests.
func TestEmailConnector_SendMail(t *testing.T) {
	if !doTests {
		return
	}

	con := mvcapp.NewEmailConnector("smtp.yourhost.com", 587, "you@domain.tld", "secret-password")
	msg, err := mvcapp.NewEmailMessage("From <You@Domain.tld>", "Recipient <recipient@domain.tld>", "Testing Email", "<strong>Test email</strong><br /><br /> This email was generated via Golang unit tests to proof ensure that the MvcApp Email Connector is working as expected.")

	if err != nil {
		t.Error(err)
	}

	err = con.SendMail(msg)
	if err != nil {
		t.Error(err)
	}
}

// TestEmailConnector_SendTemplateMail ensures that the EmailConnector.SendTemplateMail method operates
// as expected Note you will need to modify this method with your email credentials / recipients and set
// doTests to true to include this with your unit tests.
func TestEmailConnector_SendTemplateMail(t *testing.T) {
	if !doTests {
		return
	}

	con := mvcapp.NewEmailConnector("smtp.yourhost.com", 587, "you@domain.tld", "secret-password")
	data := struct {
		CompanyName string
	}{
		CompanyName: "Digivance Technologies",
	}

	tdata := "{{ define \"EmailMessage\" }}<h1>Template Test</h1>The maker of this library is {{ .CompanyName }}!<br /><br />{{ end }}"
	tpath := fmt.Sprintf("%s/testMail.htm", mvcapp.GetApplicationPath())
	ioutil.WriteFile(tpath, []byte(tdata), 0644)
	defer os.RemoveAll(tpath)

	msg, err := mvcapp.NewEmailFromTemplate("From <You@Domain.tld>", "Recipient <recipient@domain.tld>", "Template Test!", tpath, data)

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
// doTests to true to include this with your unit tests.
func TestEmailConnector_SendWithAttachment(t *testing.T) {
	if !doTests {
		return
	}

	con := mvcapp.NewEmailConnector("smtp.domain.com", 587, "username", "password")
	msg, err := mvcapp.NewEmailMessage("from", "to", "Testing Email", "<strong>Test email with CC and BCC recipients!</strong><br /><br /> This email was generated via Golang unit tests to proof ensure that the MvcApp Email Connector is working as expected.")

	if err != nil {
		t.Error(err)
	}

	if err = msg.AddAttachment(fmt.Sprintf("%s/LICENSE", mvcapp.GetApplicationPath())); err != nil {
		t.Error(err)
	}

	err = con.SendMail(msg)
	if err != nil {
		t.Error(err)
	}
}

// TestEmailConnector_SendMultiRecipients ensures that the EmailConnector.SendMultiRecipients method operates
// as expected Note you will need to modify this method with your email credentials / recipients and set
// doTests to true to include this with your unit tests.
func TestEmailConnector_SendMultiRecipients(t *testing.T) {
	if !doTests {
		return
	}

	con := mvcapp.NewEmailConnector("smtp.domain.com", 587, "username", "password")
	msg, err := mvcapp.NewEmailMessage("from", "to", "Testing Email", "<strong>Test email with CC and BCC recipients!</strong><br /><br /> This email was generated via Golang unit tests to proof ensure that the MvcApp Email Connector is working as expected.")

	if err != nil {
		t.Error(err)
	}
	if err = msg.AddRecipient("to@domain.com"); err != nil {
		t.Error(err)
	}

	if err = msg.AddCC("cc@domain.com"); err != nil {
		t.Error(err)
	}

	if err = msg.AddBCC("bcc@domain.com"); err != nil {
		t.Error(err)
	}

	err = con.SendMail(msg)
	if err != nil {
		t.Error(err)
	}
}
