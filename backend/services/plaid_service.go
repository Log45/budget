package services

import (
	"os"

	"github.com/plaid/plaid-go/v43/plaid"
)

// TODO: connect webapp to plaid
func newPlaidClient() *plaid.APIClient {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", os.Getenv("PLAID_CLIENT_ID"))
	configuration.AddDefaultHeader("PLAID-SECRET", os.Getenv("PLAID_SECRET"))
	configuration.UseEnvironment(plaid.Production)
	return plaid.NewAPIClient(configuration)
}
