package billing

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
	types "github.com/BytemarkHosting/bytemark-client/lib/billing"
)

// GetPerson gets the person object for the named user
func GetPerson(client lib.Client, username string) (person types.Person, err error) {
	req, err := client.BuildRequest("GET", lib.BillingEndpoint, "/api/v1/people?username=%s", username)
	if err != nil {
		return
	}
	people := []types.Person{}

	_, _, err = req.Run(nil, &people)

	if len(people) == 0 {
		err = fmt.Errorf("No people were returned with the username %s", username)
		return
	}
	person = people[0]
	return
}
