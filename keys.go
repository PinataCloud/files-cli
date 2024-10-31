package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func ListKeys(name string, revoked bool, limitedUse bool, exhausted bool, offset string) (KeyListResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return KeyListResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/pinata/keys?")

	params := []string{}

	if name != "" {
		params = append(params, "name="+name)
	}

	if revoked {
		params = append(params, "revoked=true")
	}

	if limitedUse {
		params = append(params, "limitedUse=true")
	}

	if exhausted {
		params = append(params, "exhausted=true")
	}
	if offset != "" {
		params = append(params, "offset="+offset)
	}

	if len(params) > 0 {
		url += strings.Join(params, "&")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return KeyListResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return KeyListResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return KeyListResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response KeyListResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return KeyListResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return KeyListResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func CreateKey(name string, admin bool, uses int, endpoints []string) (CreateKeyResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return CreateKeyResponse{}, err
	}

	payload := CreatKeyBody{
		KeyName: name,
		Permissions: Permissions{
			Admin: admin,
		},
	}

	if uses > 0 {
		payload.MaxUses = uses
	}

	if !admin && len(endpoints) > 0 {
		dataEndpoints := DataEndpoints{}
		pinningEndpoints := PinningEndpoints{}

		for _, endpoint := range endpoints {
			switch endpoint {
			case "pinList":
				dataEndpoints.PinList = true
			case "userPinnedDataTotal":
				dataEndpoints.UserPinnedDataTotal = true

			case "hashMetadata":
				pinningEndpoints.HashMetadata = true
			case "hashPinPolicy":
				pinningEndpoints.HashPinPolicy = true
			case "pinByHash":
				pinningEndpoints.PinByHash = true
			case "pinFileToIPFS":
				pinningEndpoints.PinFileToIPFS = true
			case "pinJSONToIPFS":
				pinningEndpoints.PinJSONToIPFS = true
			case "pinJobs":
				pinningEndpoints.PinJobs = true
			case "unpin":
				pinningEndpoints.Unpin = true
			case "userPinPolicy":
				pinningEndpoints.UserPinPolicy = true
			}
		}

		payload.Permissions.Endpoints = Endpoints{
			Data:    dataEndpoints,
			Pinning: pinningEndpoints,
		}
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return CreateKeyResponse{}, errors.Join(err, errors.New("Failed to marshal paylod"))
	}

	url := fmt.Sprintf("https://api.pinata.cloud/v3/pinata/keys")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return CreateKeyResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CreateKeyResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return CreateKeyResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response CreateKeyResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return CreateKeyResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return CreateKeyResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func RevokeKey(id string) error {
	jwt, err := findToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/groups/%s", id)

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

	fmt.Println("Group Deleted")

	return nil

}
