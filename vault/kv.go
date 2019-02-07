package vault

import (
	"strings"
	"log"
)

func WriteSecret(vaultToken, url, data string) {

	_, err := sendVaultRequest("POST", url, vaultToken, strings.NewReader(data))
	if err != nil {
		log.Println(err)
	}
}

func ReadSecret(vaultToken, url string) (interface{}, error){
	return sendVaultRequest("GET", url, vaultToken, nil)
}

func DeleteSecret(vaultToken, url string) {
	_, err := sendVaultRequest("DELETE", url, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
}
