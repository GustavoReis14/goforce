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
	ServerUrl         string `xml:"serverUrl"`
	SessionId         string `xml:"sessionId"`
	UserId            string `xml:"userId"`
}

type UserInfo struct {
	Username, Password, SecretToken string
}

func (c *Client) loginSoap() (*Result, error) {
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
		</env:Envelope>`, c.userInfo.Username, (c.userInfo.Password + c.userInfo.SecretToken))

	payload := strings.NewReader(rawPayload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

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
