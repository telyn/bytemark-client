package lib

import (
	"bytes"
	"net/url"
)

func (c *bytemarkClient) CreateCreditCard(cc *CreditCard) (ref string, err error) {
	req, err := c.BuildRequestNoAuth("POST", EP_SPP, "/card.ref")
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
