package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func DeleteFile(id string) error {
	jwt, err := findToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/%s", id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server Returned an error %d, check CID", resp.StatusCode)
	}

	fmt.Println("File Deleted")

	return nil

}

func ListFiles(amount string, pageToken string, cidPending bool) (ListResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return ListResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files?")

	params := []string{}

	if amount != "" {
		params = append(params, "limit="+amount)
	}

	if pageToken != "" {
		url += "&pagetoken=" + pageToken
	}

	if cidPending {
		url += "&cidPending=true"
	}

	if len(params) > 0 {
		url += strings.Join(params, "&")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ListResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ListResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ListResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response ListResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ListResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return ListResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}
