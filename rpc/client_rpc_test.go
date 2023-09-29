package rpc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestEncodePrivateTxParams(t *testing.T) {
	signedTx := "0x0"
	hints := Hints{
		CallData: true,
		Logs:     true,
	}
	options := PrivateTxOptions{
		Hints:          hints,
		MaxBlockNumber: hexutil.Uint64(100),
		Builders:       []string{"builder1", "builder2"},
	}

	expectedData := struct {
		Tx             string         `json:"tx"`
		MaxBlockNumber hexutil.Uint64 `json:"maxBlockNumber,omitempty"`
		Preferences    struct {
			Fast     bool     `json:"fast"`
			Privacy  []string `json:"privacy,omitempty"`
			Builders []string `json:"builders,omitempty"`
		} `json:"preferences"`
	}{
		Tx:             signedTx,
		MaxBlockNumber: options.MaxBlockNumber,
		Preferences: struct {
			Fast     bool     `json:"fast"`
			Privacy  []string `json:"privacy,omitempty"`
			Builders []string `json:"builders,omitempty"`
		}{
			Fast:     true,
			Privacy:  hints.String(),
			Builders: options.Builders,
		},
	}

	actualData := encodePrivateTxParams(signedTx, &options).(struct {
		Tx             string         `json:"tx"`
		MaxBlockNumber hexutil.Uint64 `json:"maxBlockNumber,omitempty"`
		Preferences    struct {
			Fast     bool     `json:"fast"`
			Privacy  []string `json:"privacy,omitempty"`
			Builders []string `json:"builders,omitempty"`
		} `json:"preferences"`
	})

	if actualData.Tx != expectedData.Tx ||
		actualData.MaxBlockNumber != expectedData.MaxBlockNumber ||
		actualData.Preferences.Fast != expectedData.Preferences.Fast ||
		len(actualData.Preferences.Privacy) != len(expectedData.Preferences.Privacy) ||
		len(actualData.Preferences.Builders) != len(expectedData.Preferences.Builders) {
		t.Errorf("encodePrivateTxParams did not produce the expected result.")
	}

	for i, actualPrivacy := range actualData.Preferences.Privacy {
		if actualPrivacy != expectedData.Preferences.Privacy[i] {
			t.Errorf("Privacy hint mismatch at index %d. Expected: %s, Actual: %s", i, expectedData.Preferences.Privacy[i], actualPrivacy)
		}
	}

	for i, actualBuilder := range actualData.Preferences.Builders {
		if actualBuilder != expectedData.Preferences.Builders[i] {
			t.Errorf("Builder mismatch at index %d. Expected: %s, Actual: %s", i, expectedData.Preferences.Builders[i], actualBuilder)
		}
	}
}

func TestHints_String(t *testing.T) {
	hints := Hints{
		CallData:         true,
		ContractAddress:  true,
		FunctionSelector: true,
		Logs:             true,
		TxHash:           true,
		Hash:             true,
		SpecialLogs:      true,
	}

	expectedHints := []string{"calldata", "contract_address", "function_selector", "logs", "tx_hash", "hash", "special_logs"}
	actualHints := hints.String()

	if len(actualHints) != len(expectedHints) {
		t.Errorf("Expected hints length: %d, but got: %d", len(expectedHints), len(actualHints))
	}

	for i, hint := range expectedHints {
		if actualHints[i] != hint {
			t.Errorf("Expected hint: %s, but got: %s", hint, actualHints[i])
		}
	}
}
