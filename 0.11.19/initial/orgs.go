package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Organization ...
type Organization struct {
	UserName    string `json:"username"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	WebSite     string `json:"website"`
	Location    string `json:"location"`
}

// NewGogsOrg ...
func NewGogsOrg(name, owner string) {
	url := gogsHost + "/api/v1/admin/users/" + owner + "/orgs"
	org := &Organization{
		UserName: name,
		FullName: "",
	}

	outJSON, err := json.Marshal(org)
	if err != nil {
		fmt.Println(err)
	}
	payload := strings.NewReader(string(outJSON))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Basic ")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
