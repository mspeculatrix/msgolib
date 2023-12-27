/*
Package emailutils
Library: msgolib
Offered up under GPL 3.0 but absolutely not guaranteed fit for use.
This is code created by an amateur dilettante, so use at your own risk.
Github: https://github.com/mspeculatrix
Blog: https://mansfield-devine.com/speculatrix/

The best way to use this is with one or more config files holding details of
the email settings. The config file should consist of:

host=mail.example.com
port=465
user=user@example.com
pass=top_secret_password
use_tls=no

If 'port' is excluded, the code will select 465 or 587 depending on the value of
'use_tls' (if 'yes', uses 465, any other value will cause 587 to be selected).
If 'use_tls' is ommitted, it will default to 'no'.

The standard location for a config file is: /etc/email/email_default.cfg

But these details can also be put manually in a map[string]string where the keys
are: host, port, user, pass and use_tls

To use this library:

Create the config map manually or by reading the details from a file with:
	emailConfig := ReadConfigFile(<filepath>)

Create an EmailMessage object - eg,

var email emailutils.EmailMessage

Set the headers:
	email.AddRecipient("recipient@example.com")
	// Could be more than one recipient
	email.AddRecipient("other@example.com")
	email.SetSender("sender@example.com")
	email.SetSenderName("Mr A Sender")
	email.SetSubject("A message just for you")

The email.Header.From setting could use the 'user' entry from the config.
In the example above, this would be emailConfig["user"].

Build the email body:

	email.BodyAppend("IoT Alert\r\n---------)
	email.BodyAppend("Some text")
	email.BodyAppend("Some more text")

Add a signature if you want.

	email.AddSignature("My sig")

Send the email:

	email.SendEmail(emailCfg)

*/

package emailutils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/mspeculatrix/msgolib/fileutils"
)

var (
	ConfigKeys = []string{"host", "port", "user", "pass", "use_tls"}
)

// type EmailError struct {
// 	msg string
// 	err error
// }

// func (e *EmailError) Error() string {
// 	return e.msg + e.err.Error()
// }

type EmailHeader struct {
	Name    string   // sender's name
	From    string   // sender's email address
	To      []string // recipients' email addresses
	Subject string   // message subject line
}

type EmailMessage struct {
	Header EmailHeader
	Body   string
}

// AddRecipient() adds an email address to the EmailHeader.To list
func (e *EmailMessage) AddRecipient(addr string) {
	e.Header.To = append(e.Header.To, addr)
}

// AddSignature() simply adds a given string to the end of the message.
// Normally, this string will be a constant in any given program.
// Usage: emailutils.AddSignature(&msg, "Sig text")
func (e *EmailMessage) AddSignature(sigStr string) {
	e.Body += "\r\n" + "--" + "\r\n" + sigStr + "\r\n"
}

// BodySet() creates the main text of the email.
func (e *EmailMessage) BodySet(text string) {
	e.Body = text + "\r\n"
}

// BodyAppend() adds text to any existing text in the main
// body of the email.
func (e *EmailMessage) BodyAppend(text string) {
	e.Body += text + "\r\n"
}

// CheckHeaders() checks to see if the essential headers have values
func (hdr EmailHeader) CheckHeaders() error {
	errs := make([]string, 0)
	if hdr.From == "" {
		errs = append(errs, "Header from field empty")
	}
	if len(hdr.To) == 0 {
		errs = append(errs, "Header to field empty")
	}
	if hdr.Subject == "" {
		errs = append(errs, "Header subject empty")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs[:], "; "))
	}
	return nil
}

// HeaderString() creates a single string from the headers.
func (e EmailMessage) HeaderString() string {
	hdrStr := ""
	fromStr := ""
	err := e.Header.CheckHeaders()
	if err != nil {
		log.Println(err)
	} else {
		if len(e.Header.Name) == 0 {
			fromStr = e.Header.From
		} else {
			fromStr = e.Header.Name + " <" + e.Header.From + ">"
		}
		hdrStr = "From: " + fromStr + "\r\n"
		hdrStr += "To: " + strings.Join(e.Header.To, ",") + "\r\n"
		hdrStr += "Subject: " + e.Header.Subject + "\r\n\r\n"
	}
	return hdrStr
}

// ReadConfigFile() reads a configuration file using my own
// fileutils.ReadConfigFile() function.
func ReadConfigFile(cfgFile string) (map[string]string, error) {
	var errs = []string{}
	var err error = nil

	emailCfg, err := fileutils.ReadConfigFile(cfgFile)
	if err != nil {
		errs = append(errs, "Error reading email config file")
	}
	// make sure all the important keys are present and had some value in the
	// config file
	for _, v := range ConfigKeys {
		if len(emailCfg[v]) == 0 {
			switch v {
			case "port":
				// no port was provided, so we'll go with default values
				emailCfg["port"] = "587"
				if emailCfg["use_tls"] == "yes" {
					emailCfg["port"] = "465"
				}
			case "use_tls":
				emailCfg["use_tls"] = "no"
			default:
				errs = append(errs, "Missing config value: "+v)
			}
		}
	}
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs[:], "; "))
	}
	return emailCfg, err
}

// SendEmail() does what it says on the tin
func (msg EmailMessage) SendEmail(cfg map[string]string) error {
	//var errors = []string{}
	var e_err error
	email := msg.HeaderString()
	email += msg.Body + "\r\n"
	err := msg.Header.CheckHeaders()
	if err != nil {
		e_err = fmt.Errorf("Error in email headers " + err.Error())
	} else {
		auth := smtp.PlainAuth("", cfg["user"], cfg["pass"], cfg["host"])
		if cfg["use_tls"] == "yes" {
			// Create a TLS connection.
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         cfg["host"],
			}
			for _, mailto := range msg.Header.To {
				fmt.Println(mailto)
				conn, err := tls.Dial("tcp", cfg["host"]+":"+cfg["port"], tlsConfig)
				if err != nil {
					e_err = fmt.Errorf("Failure to make TCP connection " + err.Error())
				} else {
					c, err := smtp.NewClient(conn, cfg["host"])
					if err != nil {
						e_err = fmt.Errorf("Error creating client " + err.Error())
					} else {
						if err = c.Auth(auth); err != nil {
							e_err = fmt.Errorf("Error authenticating " + err.Error())
						} else {
							if err = c.Mail(msg.Header.From); err != nil {
								e_err = fmt.Errorf("Error setting From address " + err.Error())
							} else {
								if err = c.Rcpt(mailto); err != nil {
									e_err = fmt.Errorf("Error setting To address " + err.Error())
								} else {
									w, err := c.Data()
									if err != nil {
										e_err = fmt.Errorf("Error creating data object " + err.Error())
									} else {
										_, err = w.Write([]byte(msg.Body))
										if err != nil {
											e_err = fmt.Errorf("Error writing message " + err.Error())
										}
										w.Close()
									}
								}
							}
						}
					}
					c.Quit()
				}
			}
		} else {
			// Use the normal SendMail() method. This will employ TLS via starttls
			// if possible. It's probably fine for most uses.
			err = smtp.SendMail(cfg["host"]+":"+cfg["port"], auth, msg.Header.From,
				msg.Header.To, []byte(email))
			if err != nil {
				e_err = fmt.Errorf("Error sending email " + err.Error())
			}
		}
	}
	return e_err
}

// SetSender() sets email address of sender
func (e *EmailMessage) SetSender(sender string) {
	e.Header.From = sender
}

// SetSenderName() sets the name string of the sender
func (e *EmailMessage) SetSenderName(name string) {
	e.Header.Name = name
}

// SetSubject() sets the subject line for the message
func (e *EmailMessage) SetSubject(subject string) {
	e.Header.Subject = subject
}
