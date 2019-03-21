package vault

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// leaseSecEng is a SecEng implementation which is meant to work for Vault secrets engines
// producing leases. As well as the changepassSecEng, it communicates with endpoint creds/rolename
// to receive creds. But calling /creds again won't invalidate the existing creds if they were
// supplied as a lease. So, the leaseSecEng implementation does an explicit revocation request at
// the end of a booking. The secrets engines type "database" matches to this implementation.
type leaseSecEng struct {
	name                 string
	createLeaseURL       string
	revokeLeaseURL       string
	storeDataURL         string
	tuneLeaseDurationURL string
}

func (secEng leaseSecEng) getName() string {
	return secEng.name
}

// startBooking for a leaseSecEng means to create a lease in Vault and stores the returned
// credentialschange the credentials inside the respective kv secret engine.
func (secEng leaseSecEng) startBooking(vaultToken, _ string, _ int) {
	data, err := json.Marshal(secEng.createLease(vaultToken))
	if err != nil {
		log.Println(err)
	}
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

// endBooking for a leaseSecEng deletes the stored credentials from kv storage and then
// revokes all leases associated with the secrets engine for the configured role.
func (secEng leaseSecEng) endBooking(vaultToken string) {
	vaultStorageDelete(vaultToken, secEng.storeDataURL)
	secEng.revokeLease(vaultToken)
}

func (secEng leaseSecEng) readCreds(vaultToken string) (map[string]interface{}, error) {
	return vaultStorageRead(vaultToken, secEng.storeDataURL)
}

func (secEng leaseSecEng) createLease(vaultToken string) map[string]interface{} {
	data, err := sendVaultDataRequest("GET", secEng.createLeaseURL, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

func (secEng leaseSecEng) revokeLease(vaultToken string) {
	err := sendVaultRequestEmtpyResponse("POST", secEng.revokeLeaseURL, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
}

func tuneLeaseDuration(tuneLeaseDurationURL string, maxBookingDays int) {
	hours := maxBookingDays * 24
	payload := fmt.Sprintf("{\"default_lease_ttl\": \"%vh\", \"max_lease_ttl\": \"%vh\"}", hours, hours)
	vaultToken := CreateVaultToken()
	err := sendVaultRequestEmtpyResponse("POST", tuneLeaseDurationURL, vaultToken, strings.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
}
