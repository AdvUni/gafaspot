package vault

import (
	"log"
	"strings"
)

func vaultStorageWrite(vaultToken, url, data string) {

	_, err := sendVaultRequest("POST", url, vaultToken, strings.NewReader(data))
	if err != nil {
		log.Println(err)
	}
}

func vaultStorageRead(vaultToken, url string) (interface{}, error) {
	return sendVaultRequest("GET", url, vaultToken, nil)
}

func vaultStorageDelete(vaultToken, url string) {
	_, err := sendVaultRequest("DELETE", url, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
}
