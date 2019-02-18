package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func sendVaultDataRequest(requestType, url, vaultToken string, body io.Reader) (interface{}, error) {
	res, err := sendVaultRequest(requestType, url, vaultToken, body)
	if err != nil {
		return nil, err
	}
	data, ok := res["data"]
	if !ok {
		err = fmt.Errorf("Malformed json response from vault: Didn't find expected field 'data'.")
		return nil, err
	}
	if data == "null" {
		err = fmt.Errorf("Json response from vault not has expected content: Tried to fetch field 'data', but it seems to be empty")
		return nil, err
	}
	return data, nil
}

func sendVaultTokenRequest(requestType, url string, body io.Reader) (string, error) {
	res, err := sendVaultRequest(requestType, url, "", body)
	if err != nil {
		return "", err
	}
	authField, ok := res["auth"]
	if !ok {
		err = fmt.Errorf("Malformed json response from vault: Didn't find expected field 'auth'.")
		return "", err
	}
	if authField == "null" {
		err = fmt.Errorf("Json response from vault not has expected content: Tried to fetch field 'auth', but it seems to be empty")
		return "", err
	}
	/*
		"auth":{
			"client_token":"s.2JlKZZKp2w5YUiiOcZdxXBVu",
			"accessor":"74CuNYNKar5XigEcAX4sTo10",
			"policies":["default"],
			"token_policies":["default"],
			"metadata":{"role_name":"testrole"}
	*/
	token, ok := authField.(map[string]interface{})["client_token"]
	if !ok {
		err = fmt.Errorf("Malformed json response from vault: Didn't find expected field 'client_token' inside 'auth'.")
		return "", err
	}
	return token.(string), nil

}

func sendVaultRequestEmtpyResponse(requestType, url, vaultToken string, body io.Reader) error {
	res, err := sendVaultRequest(requestType, url, vaultToken, body)
	if err != nil {
		return err
	}
	if res != nil {
		err = fmt.Errorf("Expected an empty response. Instead, got following content: %v", res)
		return err
	}

	return nil
}

func sendVaultRequest(requestType, url, vaultToken string, body io.Reader) (map[string]interface{}, error) {

	// Build request
	req, err := http.NewRequest(requestType, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if vaultToken != "" {
		req.Header.Set("X-Vault-Token", vaultToken)
	}
	log.Println(req)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Println(resp)

	// Parse json output to an unstructured map
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Check status code.
	// For a successful request, return the json response, otherwise the errors.
	switch status := resp.Status; status {
	case "200 OK":
		return result, nil
	case "204 No Content":
		return nil, nil
	default:
		jsonErrors := result["errors"]
		err = fmt.Errorf("Vault API returned following error(s): %v", jsonErrors)
		return nil, err
	}
}

func joinRequestPath(addressStart string, subpaths ...string) string {
	url, err := url.Parse(addressStart)
	if err != nil {
		log.Println(err)
	}

	for _, item := range subpaths {
		url.Path = path.Join(url.Path, item)
	}
	return url.String()
}

func CreateVaultToken(approleID, approleSecret, vaultAddress string) string {
	url := joinRequestPath(vaultAddress, "auth/approle/login")
	payload := fmt.Sprintf("{\"role_id\": \"%v\", \"secret_id\": \"%v\"}", approleID, approleSecret)
	token, err := sendVaultTokenRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
	return token
}
