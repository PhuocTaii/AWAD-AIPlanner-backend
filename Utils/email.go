package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(toEmail, verificationCode string) error {
	from := os.Getenv("SYSTEM_EMAIL")
	password := os.Getenv("SYSTEM_EMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	verificationLink := os.Getenv("EMAIL_VERIFICATION_URL") + verificationCode
	subject := "Subject: Email Verification\n"
	body := fmt.Sprintf("Click the link to verify your account: %s", verificationLink)
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
}
