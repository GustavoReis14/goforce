package connection

import (
	"fmt"
)

const (
	PROD_URL    = "https://login.salesforce.com"
	SANDBOX_URL = "https://test.salesforce.com"

	API_VERSION = "61.0"
)

type Client struct {
	userInfo             *UserInfo
	loginUrl, apiVersion string
	protocol             int8
	connection           *Result
}

func (c *Client) SetUserInfo(info *UserInfo) {
	c.userInfo = info
	c.apiVersion = API_VERSION
	c.loginUrl = SANDBOX_URL
	c.protocol = LOGIN_PROTOCOL_SOAP
}

func (c *Client) SetApiVersion(version string) {
	c.apiVersion = version
}

func (c *Client) SetLoginUrl(url string) {
	c.loginUrl = url
}

func (c *Client) Login(protocol int8) {
	switch protocol {
	case LOGIN_PROTOCOL_SOAP:
		res, err := c.loginSoap()
		if err != nil {
			fmt.Println("Error", err)
		} else {
			c.connection = res
			fmt.Println("Success")
			fmt.Println(*res)
		}
	}
}
