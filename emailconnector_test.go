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

func TestSendMail(t *testing.T) {
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

func TestSendTemplateMail(t *testing.T) {
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

func TestSendEmailWithAttachment(t *testing.T) {
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

func TestSendEmailMultipleRecipients(t *testing.T) {
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
