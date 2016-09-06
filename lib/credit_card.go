package lib

import (
	"bytes"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"net/url"
)

func (c *bytemarkClient) CreateCreditCard(cc *spp.CreditCard) (ref string, err error) {
	req, err := c.BuildRequestNoAuth("POST", SPPEndpoint, "/card.ref")
	if err != nil {
		return "", err
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
	_, response, err := req.Run(bytes.NewBufferString(values.Encode()), nil)

	return string(response), err
}
