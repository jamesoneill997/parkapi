package middleware

import (
	"fmt"
	"os"

	"github.com/jamesoneill997/parkapi/logs"
	"github.com/jamesoneill997/parkapi/structs"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

/*WelcomeEmail function sends an email to a user upon signing up*/
func WelcomeEmail(user structs.User) {
	//get private key that is stored as an env var
	privateKey := os.Getenv("emailkey")

	//user data
	firstName := user.FirstName
	surname := user.Surname
	toEmail := user.Email

	from := mail.NewEmail("James from Parkai", "jamesoneill997@gmail.com")
	subject := "Thank You for signing up!"

	to := mail.NewEmail(fmt.Sprintf("%s %s", firstName, surname), toEmail)
	plainTextContent := "Thanks for creating an account!"
	htmlContent := "<strong>Don't forget, you're here forever.</strong>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(privateKey)
	_, err := client.Send(message)

	if err != nil {
		logs.LogError(err)
	}

}
