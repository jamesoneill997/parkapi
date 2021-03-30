package payments

import (
	"fmt"
	"os"
	"parkapi/structs"
	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

/*CreateCustomer function creates a customer that is linked to the parkai master account*/
func CreateCustomer(user structs.User) (*stripe.Customer, error) {

	stripe.Key = os.Getenv("stripekey")

	params := &stripe.CustomerParams{
		Description: stripe.String("My First Test Customer (created for API docs)"),
		Email:       &user.Email,
	}

	return customer.New(params)
}

/*DeleteCustomer function takes a customer id and removes them as a customer from the parkai master account*/
func DeleteCustomer(custID string) (*stripe.Customer, error) {
	return customer.Del(custID, nil)
}

/*ChargeCustomer will process a charge on a customer for carpark access*/
func ChargeCustomer(custID string, amount float32) {
	stripe.Key = os.Getenv("stripekey")
	params := &stripe.ChargeParams{
		//amount should be represented in cents
		Amount:      stripe.Int64(int64(amount * 100)),
		Currency:    stripe.String(string(stripe.CurrencyEUR)),
		Description: stripe.String(fmt.Sprintf("Parking fee processed at %s", time.Now())),
		Source:      &stripe.SourceParams{Token: stripe.String("tok_visa")},
	}

	charge.New(params)
}
