package vault

import (
	"log"
	"strings"
)

func vaultStorageWrite(vaultToken, url, data string) {

	err := sendVaultRequestEmtpyResponse("POST", url, vaultToken, strings.NewReader(data))
	if err != nil {
		log.Println(err)
	}
}

func vaultStorageRead(vaultToken, url string) (interface{}, error) {
	return sendVaultDataRequest("GET", url, vaultToken, nil)
}

func vaultStorageDelete(vaultToken, url string) {
	err := sendVaultRequestEmtpyResponse("DELETE", url, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
}
