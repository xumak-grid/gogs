package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Repository represents a gogs repository
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
	Gitignores  string `json:"gitignores"`
	License     string `json:"license"`
	Readme      string `json:"readme"`
}

// NewGogsRepo create a new gogs repository
func NewGogsRepo(name, owner string) {
	url := gogsHost + "/api/v1/admin/users/" + owner + "/repos"
	repo := &Repository{
		Name:        name,
		Description: "Your repository description",
		Private:     true, // Until we resolve drone issue for resolving gogs url
		AutoInit:    true,
		Gitignores:  "Java",
		License:     "MIT License",
		Readme:      "Default",
	}

	outJSON, err := json.Marshal(repo)
	if err != nil {
		fmt.Println(err)
	}
	payload := strings.NewReader((string)(outJSON))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth("xumak", "xumakgt")
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
