package connection

const (
	PROD_URL    = "https://login.salesforce.com"
	SANDBOX_URL = "https://test.salesforce.com"

	API_VERSION = "61.0"

	LOGIN_PROTOCOL_SOAP   = 0
	LOGIN_PROTOCOL_OAUTH2 = 1

	LOGIN_PROTOCOL_SOAP_PATH   = "/services/Soap/u/"
	LOGIN_PROTOCOL_OAUTH2_PATH = "/services/oauth2/token"
)

type user struct {
	username, password, secretToken string
}

type oauth2 struct {
	clientId, clientSecret, redirectUri string
}

type Client struct {
	loginUrl, apiVersion string
	protocol             int8
	user                 user
	oauth2               oauth2

	token    string
	instance string
}

func (c *Client) SetUserInfo(userInfos ...string) {
	if len(userInfos) >= 2 {
		c.user.username = userInfos[0]
		c.user.password = userInfos[1]
	}

	if len(userInfos) >= 3 {
		c.user.secretToken = userInfos[2]
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

/*
Login method expect the following protocols:
  - LOGIN_PROTOCOL_SOAP
  - LOGIN_PROTOCOL_OAUTH2
*/
func (c *Client) Login(protocol int8) (UserInfo, error) {
	switch protocol {
	case LOGIN_PROTOCOL_SOAP:
		userInfo, err := c.loginSoap()
		if err != nil {
			return userInfo, err
		}
	case LOGIN_PROTOCOL_OAUTH2:
		err := c.loginOAuth2()
		if err != nil {
			return UserInfo{}, err
		}
	}

	return UserInfo{}, nil
}
