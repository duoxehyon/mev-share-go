package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	clientURL := "http://example.com"
	privKey, _ := crypto.GenerateKey()
	client := NewClient(clientURL, privKey)

	assert.NotNil(t, client)
	assert.Equal(t, clientURL, client.baseURL)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.privKey)
}

func TestClient_CallWithSig(t *testing.T) {
	privKey, _ := crypto.GenerateKey()
	client := NewClient("", privKey)

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.NotNil(t, r.Header.Get("X-Flashbots-Signature"))

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "success")
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client.baseURL = server.URL

	response, err := client.CallWithSig("test_method", "param1", "param2")
	assert.NoError(t, err)
	assert.Equal(t, "success", string(response))
}

func TestSimBundleOverrides_DefaultValues(t *testing.T) {
	overrides := SimBundleOverrides{}
	assert.Equal(t, hexutil.Uint64(0), overrides.ParentBlock)
	assert.Nil(t, overrides.BlockNumber)
}

func TestBundlePrivacy_MarshalJSON(t *testing.T) {
	hints := Hints{
		CallData:         true,
		ContractAddress:  true,
		FunctionSelector: true,
		Logs:             true,
		TxHash:           true,
		Hash:             true,
		SpecialLogs:      true,
	}

	bundlePrivacy := BundlePrivacy{
		Hints:    hints,
		Builders: []string{"builder1", "builder2"},
	}

	expectedJSON := `{"hints":["calldata","contract_address","function_selector","logs","tx_hash","hash","special_logs"],"builders":["builder1","builder2"]}`

	actualJSON, err := json.Marshal(bundlePrivacy)
	if err != nil {
		t.Errorf("Error marshaling BundlePrivacy: %v", err)
	}
	if string(actualJSON) != expectedJSON {
		t.Errorf("Expected JSON: %s, but got: %s", expectedJSON, string(actualJSON))
	}
}

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
