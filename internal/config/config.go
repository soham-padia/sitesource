package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/solow-crypt/bookings/internal/models"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	Infolog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
