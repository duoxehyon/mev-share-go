package rpc

import (
	"encoding/json"
	"testing"
)

func TestEncodePrivateTxParams(t *testing.T) {
	// Create a sample PrivateTxOptions struct
	options := &PrivateTxOptions{
		Hints: Hints{
			CallData: true,
			Logs:     true,
		},
		Builders:       []string{"builder1", "builder2"},
		MaxBlockNumber: 42,
	}

	signedTx := "sampleSignedTx"

	result := encodePrivateTxParams(signedTx, options)

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Error marshaling result to JSON: %v", err)
	}

	expectedJSON := `{"tx":"sampleSignedTx","maxBlockNumber":"0x2a","preferences":{"fast":true,"privacy":{"hints":["calldata","logs"]},"builders":["builder1","builder2"]}}`

	if string(resultJSON) != expectedJSON {
		t.Errorf("Result JSON does not match the expected JSON.\nExpected: %s\nActual: %s", expectedJSON, string(resultJSON))
	}
}

func TestHints_String(t *testing.T) {
	hints := Hints{
		CallData:         true,
		ContractAddress:  true,
		FunctionSelector: true,
		Logs:             true,
		TxHash:           true,
		DefaultLogs:      true,
	}

	expectedHints := []string{"calldata", "contract_address", "function_selector", "logs", "default_logs", "tx_hash"}
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
