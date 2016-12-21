package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"net/url"
)

type sppTokenRequest struct {
	Owner      *billing.Person `json:"owner,omitempty"`
	CardEnding string          `json:"card_ending"`
}

// GetSPPToken requests a token to use with SPP from bmbilling.
// If owner is nil, authenticates against bmbilling.
func (c *bytemarkClient) GetSPPToken(cc spp.CreditCard, owner *billing.Person) (token string, err error) {
	r, err := c.BuildRequestNoAuth("POST", BillingEndpoint, "/api/v1/accounts/spp_token")
	if err != nil {
		return
	}

	// i'm not really interested in whether a card number is valid, just whether it's long enough to have a last 4 digits.
	if len(cc.Number) < 4 {
		return "", errors.New("credit card number too short")
	}

	tokenRequest := sppTokenRequest{
		Owner:      owner,
		CardEnding: cc.Number[len(cc.Number)-4:],
	}
	if owner == nil {
		r.authenticate = true
	}

	js, err := json.Marshal(&tokenRequest)
	if err != nil {
		return "", err
	}

	_, res, err := r.Run(bytes.NewBuffer(js), nil)
	token = string(res)
	return
}

// CreateCreditCard creates a credit card on SPP using the given token. Tokens must be acquired by using GetSPPToken or GetSPPTokenWithAccount first.
func (c *bytemarkClient) CreateCreditCardWithToken(cc *spp.CreditCard, token string) (ref string, err error) {
	req, err := c.BuildRequestNoAuth("POST", SPPEndpoint, "/card.ref")
	req.sppToken = token
	if err != nil {
		return
	}
	values := url.Values{}
	values.Add("account_number", cc.Number)
	values.Add("name", cc.Name)
	values.Add("expiry", cc.Expiry)
	values.Add("cvv", cc.CVV)
	if cc.Street != "" {
		values.Add("street", cc.Street)
		values.Add("city", cc.City)
		values.Add("postcode", cc.Postcode)
		values.Add("country", cc.Country)
	}
	// prevent CC details and card reference being written to log
	// this is a bit of a sledgehammer
	// TODO make it not a sledgehammer somehow
	oldfile := log.LogFile
	log.LogFile = nil
	_, response, err := req.Run(bytes.NewBufferString(values.Encode()), nil)
	log.LogFile = oldfile

	return string(response), err
}

// CreateCreditCard creates a credit card on SPP. It uses GetSPPToken to get a token.
func (c *bytemarkClient) CreateCreditCard(cc *spp.CreditCard) (ref string, err error) {
	token, err := c.GetSPPToken(*cc, nil)
	if err != nil {
		return
	}
	return c.CreateCreditCardWithToken(cc, token)

}
