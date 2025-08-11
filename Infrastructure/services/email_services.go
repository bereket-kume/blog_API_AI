package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FrontendURL  string
}

func NewEmailService(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail, frontendURL string) *EmailService {
	return &EmailService{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		FromEmail:    fromEmail,
		FrontendURL:  frontendURL,
	}
}

func (es *EmailService) SendEmail(to, subject, htmlBody string) error {
	log.Printf("Attempting to send email to: %s via %s:%s", to, es.SMTPHost, es.SMTPPort)

	// Convert port string to int
	port, err := strconv.Atoi(es.SMTPPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	// Create SMTP authentication
	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	// Create email headers
	headers := make(map[string]string)
	headers["From"] = es.FromEmail
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build email message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// Send email using SMTP
	addr := fmt.Sprintf("%s:%d", es.SMTPHost, port)

	// For Brevo, we need to use TLS
	if es.SMTPHost == "smtp-relay.brevo.com" {
		log.Printf("Connecting to Brevo SMTP server with TLS: %s", addr)

		// Create TLS config
		tlsConfig := &tls.Config{
			ServerName: es.SMTPHost,
		}

		// Connect to SMTP server with TLS
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			log.Printf("Failed to connect to SMTP server: %v", err)
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		// Create SMTP client
		client, err := smtp.NewClient(conn, es.SMTPHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Close()

		// Authenticate
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}

		// Set sender
		if err := client.Mail(es.FromEmail); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		// Set recipient
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		// Send email data
		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}

		_, err = writer.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write email data: %w", err)
		}

		if err := writer.Close(); err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		log.Println("✅ Email sent to:", to)
		return nil
	} else {
		// For other SMTP servers, use standard SMTP
		err = smtp.SendMail(addr, auth, es.FromEmail, []string{to}, []byte(message))
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		log.Println("✅ Email sent to:", to)
		return nil
	}
}

func (es *EmailService) SendVerificationEmail(username, email, token string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", es.FrontendURL, token)
	subject := "Verify your email"

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 10px;">
			<h2 style="color: #333;">👋 Welcome, %s!</h2>
			<p style="color: #555;">Thanks for signing up. Please verify your email address to continue:</p>
			<a href="%s" style="display: inline-block; padding: 12px 24px; margin: 20px 0; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Verify Email</a>
			<p style="color: #777;">Or copy and paste this link into your browser:</p>
			<p style="word-break: break-all; color: #007BFF;">%s</p>
			<p style="font-size: 0.9em; color: #aaa;">This link expires in 24 hours.</p>
		</div>`, username, link, link)

	return es.SendEmail(email, subject, body)
}

func (es *EmailService) SendPasswordResetEmail(username, email, token string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", es.FrontendURL, token)
	subject := "Reset your password"

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 10px;">
			<h2 style="color: #333;">🔒 Password Reset</h2>
			<p style="color: #555;">Hello %s, we received a request to reset your password.</p>
			<a href="%s" style="display: inline-block; padding: 12px 24px; margin: 20px 0; background-color: #f44336; color: white; text-decoration: none; border-radius: 5px;">Reset Password</a>
			<p style="color: #777;">If the button doesn't work, use this link:</p>
			<p style="word-break: break-all; color: #007BFF;">%s</p>
			<p style="font-size: 0.9em; color: #aaa;">This link expires in 15 minutes. If you didn't request this, you can ignore this email.</p>
		</div>`, username, link, link)

	return es.SendEmail(email, subject, body)
}
