package connection

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	LOGIN_PROTOCOL_SOAP      = 1
	LOGIN_PROTOCOL_SOAP_PATH = "/services/Soap/u/"
)

type envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		LoginResponse LoginResponse `xml:"loginResponse"`
	}
}

type LoginResponse struct {
	Result Result `xml:"result"`
}

type Result struct {
	MetadataServerUrl string `xml:"metadataServerUrl"`
	PasswordExpired   bool   `xml:"passwordExpired"`
	Sandbox           bool   `xml:"sandbox"`
	ServerUrl         string `xml:"serverUrl"`
	SessionId         string `xml:"sessionId"`
	UserId            string `xml:"userId"`
	UserInfo          struct {
		XMLName                    xml.Name `xml:"userInfo"`
		AccessibilityMode          bool     `xml:"accessibilityMode"`
		ChatterExternal            bool     `xml:"chatterExternal"`
		CurrencySymbol             string   `xml:"currencySymbol"`
		OrgAttachmentFileSizeLimit int      `xml:"orgAttachmentFileSizeLimit"`
		OrgDefaultCurrencyIsoCode  string   `xml:"orgDefaultCurrencyIsoCode"`
		OrgDefaultCurrencyLocale   string   `xml:"orgDefaultCurrencyLocale"`
		OrgDisallowHtmlAttachments bool     `xml:"orgDisallowHtmlAttachments"`
		OrgHasPersonAccounts       bool     `xml:"orgHasPersonAccounts"`
		OrganizationId             string   `xml:"organizationId"`
		OrganizationMultiCurrency  bool     `xml:"organizationMultiCurrency"`
		OrganizationName           string   `xml:"organizationName"`
		ProfileId                  string   `xml:"profileId"`
		RoleId                     *string  `xml:"roleId,omitempty"`
		SessionSecondsValid        int      `xml:"sessionSecondsValid"`
		UserDefaultCurrencyIsoCode *string  `xml:"userDefaultCurrencyIsoCode,omitempty"`
		UserEmail                  string   `xml:"userEmail"`
		UserFullName               string   `xml:"userFullName"`
		UserId                     string   `xml:"userId"`
		UserLanguage               string   `xml:"userLanguage"`
		UserLocale                 string   `xml:"userLocale"`
		UserName                   string   `xml:"userName"`
		UserTimeZone               string   `xml:"userTimeZone"`
		UserType                   string   `xml:"userType"`
		UserUiSkin                 string   `xml:"userUiSkin"`
	}
}

func (c *Client) validateSoapInput() error {
	expectedFields := []string{}
	if c.userInfo.username == "" {
		expectedFields = append(expectedFields, "username")
	}

	if c.userInfo.password == "" {
		expectedFields = append(expectedFields, "password")
	}

	if c.userInfo.secretToken == "" {
		expectedFields = append(expectedFields, "secretToken")
	}

	if len(expectedFields) > 1 {
		return fmt.Errorf("required info are missing\nexpected values: %s", expectedFields)
	}

	return nil

}

func (c *Client) loginSoap() (*Result, error) {
	err := c.validateSoapInput()
	if err != nil {
		return nil, err
	}

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
		</env:Envelope>`, c.userInfo.username, (c.userInfo.password + c.userInfo.secretToken))

	payload := strings.NewReader(rawPayload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", (c.loginUrl + LOGIN_PROTOCOL_SOAP_PATH + c.apiVersion), payload)

	if err != nil {
		return &Result{}, err
	}
	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("SOAPAction", "login")
	req.Header.Add("Accept", "text/xml")

	res, err := client.Do(req)
	if err != nil {
		return &Result{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &Result{}, err
	}

	var envelope envelope
	xml.Unmarshal([]byte(body), &envelope)

	return &envelope.Body.LoginResponse.Result, nil
}
