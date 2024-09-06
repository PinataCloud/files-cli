package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func ListGroups(amount string, isPublic bool) (GroupListResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return GroupListResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/groups?")

	params := []string{}

	if amount != "" {
		params = append(params, "limit="+amount)
	}

	if isPublic {
		params = append(params, "isPublic=true")
	}
	if len(params) > 0 {
		url += strings.Join(params, "&")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GroupListResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GroupListResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GroupListResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response GroupListResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GroupListResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return GroupListResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}
