package rest

import (
	"fmt"
	con "goforce/connection"
	"io"
	"net/http"
	"strings"
)

const (
	queryPath = "/services/data/v"
)

func Query(c con.Client, soql string) {
	if soql == "" {
		return
	}

	client := &http.Client{}
	url := c.Instance + queryPath + c.ApiVersion + "/query?q=" + strings.ReplaceAll(soql, " ", "+")

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Token)

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
	fmt.Println(string(body))
}
