package services

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	FromEmail   string
	FrontendURL string
}

func NewEmailService(frontendURL string) *EmailService {
	return &EmailService{
		SMTPHost:    os.Getenv("BREVO_SMTP_HOST"),
		SMTPPort:    587, // from your env BREVO_SMTP_PORT
		SMTPUser:    os.Getenv("BREVO_SMTP_USERNAME"),
		SMTPPass:    os.Getenv("BREVO_SMTP_PASSWORD"),
		FromEmail:   os.Getenv("FROM_EMAIL"),
		FrontendURL: frontendURL,
	}
}

func (es *EmailService) SendEmail(to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.FromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(es.SMTPHost, es.SMTPPort, es.SMTPUser, es.SMTPPass)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("❌ Failed to send email: %v\n", err)
		return err
	}
	log.Println("✅ Email sent to:", to)
	return nil
}

func (es *EmailService) SendVerificationEmail(username, email, token string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", es.FrontendURL, token)
	subject := "Verify your email"

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 10px;">
			<h2>👋 Welcome, %s!</h2>
			<p>Please verify your email address:</p>
			<a href="%s" style="display:inline-block;padding:12px 24px;background-color:#4CAF50;color:white;text-decoration:none;border-radius:5px;">Verify Email</a>
			<p>If the button doesn't work, use this link:</p>
			<p>%s</p>
		</div>`, username, link, link)

	return es.SendEmail(email, subject, body)
}

func (es *EmailService) SendPasswordResetEmail(username, email, token string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", es.FrontendURL, token)
	subject := "Reset your password"

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 10px;">
			<h2>🔒 Password Reset</h2>
			<p>Hello %s, click the link below to reset your password:</p>
			<a href="%s" style="display:inline-block;padding:12px 24px;background-color:#f44336;color:white;text-decoration:none;border-radius:5px;">Reset Password</a>
			<p>If the button doesn't work, use this link:</p>
			<p>%s</p>
		</div>`, username, link, link)

	return es.SendEmail(email, subject, body)
}
