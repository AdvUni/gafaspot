package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

const (
	operateBasicPath = "operate"
	storeBasicPath   = "store"
)

type SecretEngine interface {
	StartBooking(vaultToken, sshKey string)
	EndBooking(vaultToken, sshKey string)
}

func NewSecretEngine(engineType, vaultAddress, env, name, role string) SecretEngine {
	switch engineType {
	case "ad", "ontap":
		log.Println("adding a creds secret engine")

		changeCredsURL := joinRequestPath(vaultAddress, operateBasicPath, env, name, ontapCredsPath, role)
		log.Println("creds path: ", changeCredsURL)
		storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, env, name, role, "data")
		log.Println("kv path: ", storeDataURL)

		return OntapSecretEngine{
			changeCredsURL,
			storeDataURL,
		}
	case "ssh":
		log.Println("adding ssh secret engine")

		signURL := joinRequestPath(vaultAddress, operateBasicPath, env, name, sshSignPath, role)
		log.Println("sign path: ", signURL)
		storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, env, name, role, "signature")
		log.Println("kv path: ", storeDataURL)

		return SshSecretEngine{
			signURL,
			storeDataURL,
		}

	default:
		log.Println(fmt.Errorf("Unsupported Secret Engine type: %v", engineType))
		return nil
	}
}

func sendVaultRequest(requestType, url, vaultToken string, body io.Reader) (interface{}, error) {

	// Build request
	req, err := http.NewRequest(requestType, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vaultToken)
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
	// For a successful request, return the "data" field from json output, otherwise the errors.
	switch status := resp.Status; status {
	case "200 OK":
		return result["data"], nil
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
