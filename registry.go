package kis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Registry represents the Registry information returned by the Kubota API.
type Registry struct {
	SubscriptionID    string     `json:"SubscriptionID"`
	MachineUUID       string     `json:"MachineUUID"`
	ServiceLevel      string     `json:"ServiceLevel"`
	SubscriptionStart string     `json:"SubscriptionStart"`
	SubscriptionEnd   string     `json:"SubscriptionEnd"`
	Timestamp         CustomTime `json:"Timestamp"`
	CreateTime        CustomTime `json:"CreateTime"`
	UpdateTime        CustomTime `json:"UpdateTime"`
}

// GetRegistryByMobilePhone retrieves registry information by mobile phone number.
func (k *kubota) GetRegistryByMobilePhone(mobilePhone string, subscription string) (Registry, error) {
	return k.getRegistry("mobilePhone", mobilePhone, subscription)
}

// GetRegistryByUserName retrieves registry information by username.
func (k *kubota) GetRegistryByUserName(userName string, subscription string) (Registry, error) {
	return k.getRegistry("userName", userName, subscription)
}

// GetRegistryByMachineUUID retrieves registry information by machine UUID.
func (k *kubota) GetRegistryByMachineUUID(machineUUID string, subscription string) (Registry, error) {
	return k.getRegistry("machineUUID", machineUUID, subscription)
}

// getRegistry is a helper function to retrieve registry information based on a given field.
func (k *kubota) getRegistry(field, value, subscription string) (Registry, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/registry?%s=%s", k.authentication.Endpoint, field, value)
	if subscription != "" {
		apiURL += "&subscription=" + subscription
	}

	// Make the request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return Registry{}, fmt.Errorf("error creating registry request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Registry{}, fmt.Errorf("error making registry request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Registry{}, fmt.Errorf("error getting registry: %s", string(body))
	}

	// Unmarshal the response
	var registryResponse struct {
		Status   int      `json:"Status"`
		Resource string   `json:"Resource"`
		Payload  Registry `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&registryResponse); err != nil {
		return Registry{}, fmt.Errorf("error decoding registry response: %w", err)
	}

	// Return the registry information
	return registryResponse.Payload, nil
}
