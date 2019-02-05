package vault

import (
	"io/ioutil"
	"log"
)

const (
	credsPath = "creds"
)

type OntapSecretEngine struct {
	VaultAddress string
	VaultPath    string
	Role         string
}

func (ontap OntapSecretEngine) ChangeCreds(vaultToken string) string {

	requestPath := joinRequestPath(ontap.VaultAddress, ontap.VaultPath, credsPath, ontap.Role)

	log.Println("repuestPath: ", requestPath)

	resp, err := sendVaultRequest("GET", requestPath, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	if resp == nil {
		log.Println("response is nill")
	}

	responseData, _ := ioutil.ReadAll(resp.Body)
	return string(responseData)

}
