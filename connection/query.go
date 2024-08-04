package connection

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) Query(soql string) {
	if soql == "" {
		return
	}

	client := &http.Client{}
	url := c.instance + qUERY_PATH + c.apiVersion + "/query?q=" + strings.ReplaceAll(soql, " ", "+")

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(body)
}
