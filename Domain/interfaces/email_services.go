package interfaces

type EmailService interface {
	SendEmail(to string, subject string, message string) error
	SendVerificationEmail(username, email, token string) error
	SendPasswordResetEmail(username, email, token string) error
}
