package agin

import (
	"github.com/go-gomail/gomail"
)

type Email struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`

	DefaultTo string `mapstructure:"defaultTo" json:"defaultTo" yaml:"defaultTo"`

	dialer *gomail.Dialer
}

func (e *Email) Init() {
	if e.Port == 0 || e.Username == "" || e.Host == "" || e.Password == "" {
		panic("init email err")
	}
	e.dialer = gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
}

func (e Email) SendToDefault(msg string, subject string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Username)
	m.SetHeader("To", e.DefaultTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	err = e.dialer.DialAndSend(m)
	return
}
