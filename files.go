package main

import (
	"bytes"
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

func GetFile(id string) (GetFileResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return GetFileResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/%s", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GetFileResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetFileResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GetFileResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}
	var response GetFileResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GetFileResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return GetFileResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func UpdateFile(id string, name string) (GetFileResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return GetFileResponse{}, err
	}
	payload := FileUpdateBody{
		Name: name,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return GetFileResponse{}, errors.Join(err, errors.New("Failed to marshal paylod"))
	}

	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/%s", id)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return GetFileResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetFileResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GetFileResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}
	var response GetFileResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GetFileResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return GetFileResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func ListFiles(amount string, pageToken string, cidPending bool, name string, cid string, group string, mime_type string, keyvalues map[string]string) (ListResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return ListResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files?")

	params := []string{}

	if name != "" {
		params = append(params, "name="+name)
	}

	if cid != "" {
		params = append(params, "cid="+cid)
	}

	if group != "" {
		params = append(params, "group="+group)
	}

	if mime_type != "" {
		params = append(params, "mimeType="+mime_type)
	}

	if amount != "" {
		params = append(params, "limit="+amount)
	}

	if pageToken != "" {
		params = append(params, "pageToken="+pageToken)
	}

	if cidPending {
		params = append(params, "cidPending=true")
	}

	if len(keyvalues) > 0 {
		for key, value := range keyvalues {
			params = append(params, fmt.Sprintf("metadata[%s]=%s", key, value))
		}
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

func GetSwapHistory(cid string, domain string) (GetSwapHistoryResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return GetSwapHistoryResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/swap/%s?", cid)

	params := []string{}

	if domain != "" {
		params = append(params, "domain="+domain)
	} else {
		internalDomain, err := findGatewayDomain()
		if err != nil {
			return GetSwapHistoryResponse{}, err
		}
		params = append(params, "domain="+string(internalDomain))
	}
	if len(params) > 0 {
		url += strings.Join(params, "&")
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return GetSwapHistoryResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetSwapHistoryResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GetSwapHistoryResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response GetSwapHistoryResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GetSwapHistoryResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return GetSwapHistoryResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func AddSwap(cid string, swapCid string) (AddSwapResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return AddSwapResponse{}, err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/swap/%s", cid)

	payload := AddSwapBody{
		SwapCid: swapCid,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return AddSwapResponse{}, errors.Join(err, errors.New("Failed to marshal paylod"))
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return AddSwapResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AddSwapResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return AddSwapResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	var response AddSwapResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return AddSwapResponse{}, err
	}
	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return AddSwapResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil

}

func RemoveSwap(cid string) error {
	jwt, err := findToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.pinata.cloud/v3/files/swap/%s", cid)

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
		return fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	fmt.Println("Swap deleted")

	return nil

}
