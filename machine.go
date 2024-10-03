package kis

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
)

// Machine represents the Machine information returned by the Kubota API.
type Machine struct {
	MachineUUID        string     `json:"MachineUUID"`
	CompanyID          string     `json:"CompanyID"`
	MachineName        string     `json:"MachineName"`
	FleetID            string     `json:"FleetID"`
	EquipmentID        string     `json:"EquipmentID"`
	Brand              string     `json:"Brand"`
	Model              string     `json:"Model"`
	Type               string     `json:"Type"`
	DeviceSerialNumber string     `json:"DeviceSerialNumber"`
	SubscriptionEnd    string     `json:"SubscriptionEnd"`
	Timestamp          CustomTime `json:"Timestamp"`
	CreateTime         CustomTime `json:"CreateTime"`
	UpdateTime         CustomTime `json:"UpdateTime"`
}

// GetMachineByMobilePhone retrieves machine information by mobile phone number.
func (k *Kubota) GetMachineByMobilePhone(mobilePhone string, subscription string) (Machine, error) {
	return k.getMachine("mobilePhone", mobilePhone, subscription)
}

// GetMachineByUserName retrieves machine information by username.
func (k *Kubota) GetMachineByUserName(userName string, subscription string) (Machine, error) {
	return k.getMachine("userName", userName, subscription)
}

// GetMachineByMachineUUID retrieves machine information by machine UUID.
func (k *Kubota) GetMachineByMachineUUID(machineUUID string, subscription string) (Machine, error) {
	return k.getMachine("machineUUID", machineUUID, subscription)
}

// getMachine is a helper function to retrieve machine information based on a given field.
func (k *Kubota) getMachine(field, value, subscription string) (Machine, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/machine?%s=%s", k.authentication.Endpoint, field, value)
	if subscription != "" {
		apiURL += "&subscription=" + subscription
	}

	// Make the request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return Machine{}, fmt.Errorf("error creating machine request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Machine{}, fmt.Errorf("error making machine request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Machine{}, fmt.Errorf("error getting token: %s with statuscode: %s", string(body), string(resp.StatusCode))
	}

	// Unmarshal the response
	var machineResponse struct {
		Status   int     `json:"Status"`
		Resource string  `json:"Resource"`
		Payload  Machine `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&machineResponse); err != nil {
		return Machine{}, fmt.Errorf("error decoding machine response: %w", err)
	}

	// Return the machine information
	return machineResponse.Payload, nil
}
