package lib

import (
	"bytes"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"github.com/BytemarkHosting/bytemark-client/util/log"
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
	// prevent CC details and card reference being written to log
	// this is a bit of a sledgehammer
	// TODO make it not a sledgehammer somehow
	oldfile := log.LogFile
	log.LogFile = nil
	_, response, err := req.Run(bytes.NewBufferString(values.Encode()), nil)
	log.LogFile = oldfile

	return string(response), err
}
