package vault

import (
	"io/ioutil"
	"path"
	"log"
)

type OntapSecretEngine struct {
	VaultAddress string
	VaultPath    string
	Role         string
}

func (ontap OntapSecretEngine) ChangeCreds(vaultToken string) string {
	requestPath := path.Join(ontap.VaultAddress, ontap.VaultPath, ontap.Role)

	resp, err := sendVaultRequest("GET", requestPath, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	if resp == nil {
		log.Println("response is nill")
	}
	defer resp.Body.Close()

	responseData, _ := ioutil.ReadAll(resp.Body)
	return string(responseData)

}
