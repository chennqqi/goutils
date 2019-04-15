package logrus

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// MailHook to send logs via syslog.
type MailHook struct {
	conf   *MailConf
	levels []logrus.Level
	mailCh chan *gomail.Message
}

type MailConf struct {
	Addr string   `json:"addr" yaml:"addr"`
	Smtp string   `json:"smtp" yaml:"smtp"`
	User string   `json:"user" yaml:"user"`
	Pass string   `json:"pass" yaml:"pass"`
	Port int      `json:"port" yaml:"port"`
	To   []string `json:"to" yaml:"to"`
}

func NewMailHook(conf *MailConf, levels []logrus.Level) (*MailHook, error) {
	var h = &MailHook{conf, levels, make(chan *gomail.Message)}
	go h.run()
	return h, nil
}

func (hook *MailHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	mailCh := hook.mailCh
	conf := hook.conf

	sendMail := func(body string) {
		m := gomail.NewMessage()
		m.SetHeader("From", conf.Addr)
		m.SetHeader("To", conf.To...)
		//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "APPMON 监控报警")
		m.SetBody("text/html", body)
		//m.Attach("/home/Alex/lolcat.jpg")
		mailCh <- m
	}

	sendMail(line)
	return nil
}

func (hook *MailHook) run() {
	mailCh := hook.mailCh
	conf := hook.conf
	d := gomail.NewDialer(conf.Smtp, conf.Port, conf.User, conf.Pass)

	var s gomail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-mailCh:
			if !ok {
				return
			}
			if !open {
				if s, err = d.Dial(); err != nil {
					panic(err)
				}
				open = true
			}
			if err := gomail.Send(s, m); err != nil {
				logrus.Error("[main.sendmail] error ", err)
			}

		// Close the connection to the SMTP server if no email was sent in
		// the last 30 seconds.
		case <-time.After(30 * time.Second):
			if open {
				if err := s.Close(); err != nil {
					panic(err)
				}
				open = false
			}
		}
	}
}

func (hook *MailHook) Levels() []logrus.Level {
	return hook.levels
}
