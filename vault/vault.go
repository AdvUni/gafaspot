package vault

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

type SecretEngine interface {
	ChangeCreds(vaultToken string) string
}

func sendVaultRequest(requestType, url, vaultToken string, body io.Reader) interface{} {

	req, err := http.NewRequest(requestType, url, body)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vaultToken)
	log.Println(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	log.Println(resp)

	//extract field "data" from json output
	var result map[string]interface{}
	// TODO: catch error
	json.NewDecoder(resp.Body).Decode(&result)

	//TODO: check ecit code
	return result["data"]
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
