package connection

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	PROD_URL    = "https://login.salesforce.com"
	SANDBOX_URL = "https://test.salesforce.com"

	LOGIN_PROTOCOL_SOAP      = 1
	LOGIN_PROTOCOL_SOAP_PATH = "/services/Soap/u/"
)

type UserInfo struct {
	Username, Password, SecretToken string
}

type Client struct {
	userInfo             UserInfo
	loginUrl, apiVersion string
	protocol             int8
}

func (c *Client) SetUserInfo(info UserInfo) {
	c.userInfo = info
	c.apiVersion = "61.0"
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
		c.loginSoap()
	}
}

func (c *Client) loginSoap() {
	url := c.loginUrl + LOGIN_PROTOCOL_SOAP_PATH + c.apiVersion
	method := "POST"

	rawPayload := fmt.Sprintf(
		`<?xml version="1.0" encoding="utf-8" ?>
		<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema"
			xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
			xmlns:env="http://schemas.xmlsoap.org/soap/envelope/">
		<env:Body>
			<n1:login xmlns:n1="urn:partner.soap.sforce.com">
			<n1:username><![CDATA[%s]]></n1:username>
			<n1:password><![CDATA[%s]]></n1:password>
			</n1:login>
		</env:Body>
		</env:Envelope>`, c.userInfo.Username, (c.userInfo.Password))

	fmt.Println(rawPayload)
	payload := strings.NewReader(rawPayload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("SOAPAction", "login")
	req.Header.Add("Accept", "text/xml")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
