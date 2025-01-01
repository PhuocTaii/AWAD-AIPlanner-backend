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

	verificationLink := os.Getenv("BACKEND_URL") + "/api/auth/verify?code=" + verificationCode
	subject := "Subject: Email Verification\n"
	body := fmt.Sprintf("Click the link to verify your account: %s", verificationLink)
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
}

func SendResetPasswordEmail(toEmail, newPassword string) error {
	from := os.Getenv("SYSTEM_EMAIL")
	password := os.Getenv("SYSTEM_EMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: Reset Password\n"
	body := fmt.Sprintf("Here is your new password, please change it after login: %s", newPassword)
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
}
