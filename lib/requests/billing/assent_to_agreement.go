package billing

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	types "github.com/BytemarkHosting/bytemark-client/lib/billing"
)

func AssentToAgreement(client lib.Client, assent types.Assent) (err error) {
	req, err := client.BuildRequest("POST", lib.BillingEndpoint, "/api/v1/agreements/%s/assents", assent.AgreementID)
	if err != nil {
		return
	}
	_, _, err = req.MarshalAndRun(map[string]interface{}{
		"account_id": assent.AccountID,
		"person_id":  assent.PersonID,
		"name":       assent.Name,
	}, nil)
	return
}
