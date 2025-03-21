package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetGroup(id string) (GroupCreateResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return GroupCreateResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/groups/%s", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GroupCreateResponse{}, fmt.Errorf("server Returned an error %d, check CID", resp.StatusCode)
	}
	var response GroupCreateResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GroupCreateResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return GroupCreateResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func ListGroups(amount string, isPublic bool, name string, token string) (GroupListResponse, error) {
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

	if name != "" {
		params = append(params, "name="+name)
	}

	if token != "" {
		params = append(params, "pageToken="+token)
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

func CreateGroup(name string, isPublic bool) (GroupCreateResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return GroupCreateResponse{}, err
	}

	payload := GroupCreateBody{
		Name:     name,
		IsPublic: isPublic,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("Failed to marshal paylod"))
	}

	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/groups")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GroupCreateResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response GroupCreateResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GroupCreateResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return GroupCreateResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func UpdateGroup(id string, name string, isPublic bool) (GroupCreateResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return GroupCreateResponse{}, err
	}

	payload := GroupCreateBody{
		Name:     name,
		IsPublic: isPublic,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("Failed to marshal paylod"))
	}

	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/groups/%s", id)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GroupCreateResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GroupCreateResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response GroupCreateResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GroupCreateResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return GroupCreateResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func DeleteGroup(id string) error {
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

func AddFile(groupId string, fileId string) error {

	jwt, err := findToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/groups/%s/ids/%s", groupId, fileId)

	req, err := http.NewRequest("PUT", url, nil)
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

	fmt.Println("File added to group")

	return nil
}

func RemoveFile(groupId string, fileId string) error {

	jwt, err := findToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/groups/%s/ids/%s", groupId, fileId)

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

	fmt.Println("File removed from group")

	return nil
}
