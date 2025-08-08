package services

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	Client      *resend.Client
	FromEmail   string
	FrontendURL string
}

func NewEmailService(apiKey, fromEmail, frontendURL string) *EmailService {
	client := resend.NewClient(apiKey)
	return &EmailService{
		Client:      client,
		FromEmail:   fromEmail,
		FrontendURL: frontendURL,
	}
}

func (es *EmailService) SendEmail(to, subject, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := es.Client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send email: %v\n", err)
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
			<p style="font-size: 0.9em; color: #aaa;">This link expires in 15 minutes. If you didn’t request this, you can ignore this email.</p>
		</div>`, username, link, link)

	return es.SendEmail(email, subject, body)
}
