package email

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
	"os"
)

// SendgridService implements the Service interface using SendGrid
type SendgridService struct {
	apiKey      string
	senderName  string
	senderEmail string
}

// NewSendgridService returns a new instance of SendgridService
func NewSendgridService() Service {
	s := SendgridService{
		apiKey:      os.Getenv("SENDGRID_API_KEY"),
		senderName:  os.Getenv("SENDGRID_SENDER_NAME"),
		senderEmail: os.Getenv("SENDGRID_SENDER_EMAIL"),
	}
	var service Service = &s
	return service
}

// Send sends an email using SendGrid
func (s *SendgridService) Send(emailTo, subject, plainTextContent, htmlContent string) error {
	from := mail.NewEmail(s.senderName, s.senderEmail)
	to := mail.NewEmail(emailTo, emailTo)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	if err != nil {
		zap.L().Error("Sendgrid email send", zap.Error(err))
		return err
	}

	return nil
}
