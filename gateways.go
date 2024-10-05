package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func findGatewayDomain() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dotFilePath := filepath.Join(homeDir, ".pinata-files-cli-gateway")
	Domain, err := os.ReadFile(dotFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("JWT not found. Please authorize first using the 'auth' command")
		} else {
			return nil, err
		}
	}
	return Domain, err
}

func SetGateway(domain string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	p := filepath.Join(home, ".pinata-files-cli-gateway")
	err = os.WriteFile(p, []byte(domain), 0600)
	if err != nil {
		return err
	}

	fmt.Println("Gateway Saved!")
	return nil
}

func GetSignedURL(cid string, expires int) (GetSignedURLResponse, error) {

	jwt, err := findToken()
	if err != nil {
		return GetSignedURLResponse{}, err
	}

	domain, err := findGatewayDomain()
	if err != nil {
		return GetSignedURLResponse{}, err
	}

	domainUrl := fmt.Sprintf("https://%s/files/%s", domain, cid)

	currentTime := time.Now().Unix()

	payload := GetSignedURLBody{
		URL:     domainUrl,
		Expires: expires,
		Date:    currentTime,
		Method:  "GET",
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return GetSignedURLResponse{}, errors.Join(err, errors.New("Failed to marshal paylod"))
	}

	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/sign")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return GetSignedURLResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetSignedURLResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GetSignedURLResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response GetSignedURLResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GetSignedURLResponse{}, err
	}

	// Unescape the URL
	unescapedURL := strings.ReplaceAll(response.Data, "\\u0026", "&")
	unescapedURL = strings.Trim(unescapedURL, "\"")

	fmt.Println(unescapedURL)

	return response, nil
}
