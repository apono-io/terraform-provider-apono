package testcommon

import (
	"testing"
)

func GetTestConnectorID(t *testing.T) string {
	if IsTestAccount(t) {
		return "terraofrm-tests-account-connector"
	}
	connector, err := GetFirstConnectorV3(t)
	if err != nil {
		t.Fatalf("failed to get connector: %v", err)
	}
	return connector.ID
}
