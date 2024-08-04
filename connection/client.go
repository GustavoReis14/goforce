package connection

const (
	PROD_URL    = "https://login.salesforce.com"
	SANDBOX_URL = "https://test.salesforce.com"

	API_VERSION = "61.0"
)

type userInfo struct {
	username, password, secretToken string
}

type oauth2 struct {
	clientId, clientSecret, redirectUri string
}

type Client struct {
	loginUrl, apiVersion string
	protocol             int8
	connection           *Result
	userInfo             userInfo
	oauth2               oauth2
}

func (c *Client) SetUserInfo(userInfos ...string) {
	if len(userInfos) >= 2 {
		c.userInfo.username = userInfos[0]
		c.userInfo.password = userInfos[1]
	}

	if len(userInfos) >= 3 {
		c.userInfo.secretToken = userInfos[2]
	}

	c.apiVersion = API_VERSION
	c.loginUrl = PROD_URL
	c.protocol = LOGIN_PROTOCOL_SOAP
}

func (c *Client) SetClientCredentials(oauth2 ...string) {
	if len(oauth2) >= 2 {
		c.oauth2.clientId = oauth2[0]
		c.oauth2.clientSecret = oauth2[1]
	}

	if len(oauth2) >= 3 {
		c.oauth2.redirectUri = oauth2[2]
	}
}

func (c *Client) SetApiVersion(version string) {
	c.apiVersion = version
}

func (c *Client) SetLoginUrl(url string) {
	c.loginUrl = url
}

func (c *Client) Login(protocol int8) error {
	switch protocol {
	case LOGIN_PROTOCOL_SOAP:
		res, err := c.loginSoap()
		if err != nil {
			return err
		} else {
			c.connection = res
		}
	}

	return nil
}
