package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type oauth2Response struct {
	Token       string `json:"access_token"`
	InstanceUrl string `json:"instance_url"`
	Id          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

func (c *Client) validateOAuthInput() error {
	expectedFields := []string{}
	if c.user.username == "" {
		expectedFields = append(expectedFields, "username")
	}

	if c.user.password == "" {
		expectedFields = append(expectedFields, "password")
	}

	if c.oauth2.clientId == "" {
		expectedFields = append(expectedFields, "ClientId")
	}

	if c.oauth2.clientSecret == "" {
		expectedFields = append(expectedFields, "ClientSecret")
	}

	if len(expectedFields) > 1 {
		return fmt.Errorf("required info are missing\nexpected values: %s", expectedFields)
	}

	return nil
}

func (c *Client) loginOAuth2() error {
	err := c.validateOAuthInput()
	if err != nil {
		return err
	}

	rawPayload := fmt.Sprintf(
		"grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
		c.oauth2.clientId, c.oauth2.clientSecret, c.user.username, (c.user.password + c.user.secretToken))
	payload := strings.NewReader(rawPayload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", (c.loginUrl + lOGIN_PROTOCOL_OAUTH2_PATH), payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	response := oauth2Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error deserializing JSON: %v\n", err)
		return nil
	}

	c.token = response.TokenType + " " + response.Token
	c.instance = response.InstanceUrl

	return nil
}
