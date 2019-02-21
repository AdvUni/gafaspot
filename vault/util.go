package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

// ErrAuth is thrown if an authentication against LDAP over Vault fails for any reason.
var ErrAuth = errors.New("ldap authentication failed")

func sendVaultDataRequest(requestType, url, vaultToken string, body io.Reader) (interface{}, error) {
	res, err := sendVaultRequest(requestType, url, vaultToken, body)
	if err != nil {
		return nil, err
	}
	data, ok := res["data"]
	if !ok {
		err = fmt.Errorf("malformed json response from vault: Didn't find expected field 'data'")
		return nil, err
	}
	if data == "null" {
		err = fmt.Errorf("json response from vault not has expected content: Tried to fetch field 'data', but it seems to be empty")
		return nil, err
	}
	return data, nil
}

func sendVaultTokenRequest(url string, body io.Reader) (string, error) {
	res, err := sendVaultRequest("POST", url, "", body)
	if err != nil {
		return "", err
	}
	authField, ok := res["auth"]
	if !ok {
		err = fmt.Errorf("malformed json response from vault: Didn't find expected field 'auth'")
		return "", err
	}
	if authField == "null" {
		err = fmt.Errorf("json response from vault not has expected content: Tried to fetch field 'auth', but it seems to be empty")
		return "", err
	}
	token, ok := authField.(map[string]interface{})["client_token"]
	if !ok {
		err = fmt.Errorf("malformed json response from vault: Didn't find expected field 'client_token' inside 'auth'")
		return "", err
	}
	return token.(string), nil

}

func sendVaultLdapRequest(url string, body io.Reader) ([]interface{}, error) {
	res, err := sendVaultRequest("POST", url, "", body)
	if err != nil {
		return nil, ErrAuth
	}
	authField, ok := res["auth"]
	if !ok {
		err = fmt.Errorf("malformed json response from vault: Didn't find expected field 'auth'")
		return nil, err
	}
	if authField == "null" {
		err = fmt.Errorf("json response from vault not has expected content: Tried to fetch field 'auth', but it seems to be empty")
		return nil, err
	}
	policies, ok := authField.(map[string]interface{})["token_policies"]
	if !ok {
		err = fmt.Errorf("malformed json response from vault: Didn't find expected field 'token_policies' inside 'auth'")
		return nil, err
	}
	if policies == "null" {
		err = fmt.Errorf("json response from vault not has expected content: Tried to fetch field 'token_policies', but it seems to be empty")
		return nil, err
	}
	policySlice, ok := policies.([]interface{})
	if !ok {
		err = fmt.Errorf("malformed json response from vault: Expected field 'token_policies' to contain a list, but it doesn't. Content is instead: %v", policies)
		return nil, err
	}
	return policySlice, nil
}

func sendVaultRequestEmtpyResponse(requestType, url, vaultToken string, body io.Reader) error {
	res, err := sendVaultRequest(requestType, url, vaultToken, body)
	if err != nil {
		return err
	}
	if res != nil {
		err = fmt.Errorf("expected an empty response. Instead, got following content: %v", res)
		return err
	}

	return nil
}

// sendVaultRequest is a function to send a HTTP request of any request type towards Vault. Most request to
// Vault need a vault token as authentication. If a request is unauthenticated (such as auth requests),
// submit an empty vaultToken. As this function returnes all of Vault's json response in an interface{} map,
// there are several wrapper functions provided which unpack different values from the answer.
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
	// ignore EOF errors, as they are expected for status code 204
	if err != nil && err != io.EOF {
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
		err = fmt.Errorf("vault API returned following error(s): %v", jsonErrors)
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
