package vault

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

type SecretEngine interface {
	ChangeCreds(vaultToken string) string
}

func sendVaultRequest(requestType, url, vaultToken string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(requestType, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vaultToken)
	log.Println(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Println(resp.Header)
	log.Println(resp.Body)

	return resp, nil
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
