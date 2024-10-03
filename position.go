package kis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)


// Position represents the Position information returned by the Kubota API.
type Position struct {
	MachineUUID string     `json:"MachineUUID"`
	StatusName  string     `json:"StatusName,omitempty"`
	Latitude    float64    `json:"Latitude"`
	Longitude   float64    `json:"Longitude"`
	Speed       *float64   `json:"Speed,omitempty"` // Use a pointer to allow for null values
	Timestamp   CustomTime `json:"Timestamp"`
	CreateTime  CustomTime `json:"CreateTime"`
}

// GetLastPositionByMobilePhone retrieves the last position information by mobile phone number.
func (k *kubota) GetLastPositionByMobilePhone(mobilePhone string, subscription string) (*Position, error) {
	return k.getPosition("mobilePhone", mobilePhone, subscription)
}

// GetLastPositionByUserName retrieves the last position information by username.
func (k *kubota) GetLastPositionByUserName(userName string, subscription string) (*Position, error) {
	return k.getPosition("userName", userName, subscription)
}

// GetLastPositionByMachineUUID retrieves the last position information by machine UUID.
func (k *kubota) GetLastPositionByMachineUUID(machineUUID string, subscription string) (*Position, error) {
	return k.getPosition("machineUUID", machineUUID, subscription)
}

// GetHistoricalPositionByMobilePhone retrieves historical position information by mobile phone number.
func (k *kubota) GetHistoricalPositionByMobilePhone(mobilePhone, subscription string, startDate, endDate time.Time) ([]Position, error) {
	return k.getPositions("mobilePhone", mobilePhone, subscription, startDate, endDate)
}

// GetHistoricalPositionByUserName retrieves historical position information by username.
func (k *kubota) GetHistoricalPositionByUserName(userName, subscription string, startDate, endDate time.Time) ([]Position, error) {
	return k.getPositions("userName", userName, subscription, startDate, endDate)
}

// GetHistoricalPositionByMachineUUID retrieves historical position information by machine UUID.
func (k *kubota) GetHistoricalPositionByMachineUUID(machineUUID, subscription string, startDate, endDate time.Time) ([]Position, error) {
	return k.getPositions("machineUUID", machineUUID, subscription, startDate, endDate)
}

// getPosition is a helper function to retrieve position information based on a given field.
func (k *kubota) getPosition(field, value, subscription string) (*Position, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/position?%s=%s", k.authentication.Endpoint, field, value)
	if subscription != "" {
		apiURL += "&subscription=" + subscription
	}

	// Make the request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return &Position{}, fmt.Errorf("error creating position request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &Position{}, fmt.Errorf("error making position request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return &Position{}, fmt.Errorf("error getting position: %s", string(body))
	}

	// Unmarshal the response
	var positionResponse struct {
		Status   int       `json:"Status"`
		Resource string    `json:"Resource"`
		Payload  *Position `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&positionResponse); err != nil {
		return &Position{}, fmt.Errorf("error decoding position response: %w", err)
	}

	// Return the position information
	return positionResponse.Payload, nil
}

// getPosition is a helper function to retrieve position information based on a given field.
func (k *kubota) getPositions(field, value, subscription string, startDate, endDate time.Time) ([]Position, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/position?%s=%s", k.authentication.Endpoint, field, value)
	if subscription != "" {
		apiURL += "&subscription=" + subscription
	}
	if !startDate.IsZero() {
		s := CustomTime{startDate}
		ss, err := s.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshaling start date: %w", err)
		}
		fmt.Println(string(ss))
		apiURL += "&startDate=" + string(ss)
	}
	if !endDate.IsZero() {
		e := CustomTime{endDate}
		ee, err := e.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshaling end date: %w", err)
		}
		apiURL += "&endDate=" + string(ee)
	}
	// Make the request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating position request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making position request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("error getting position: %s", string(body))
	}
	// Unmarshal the response
	var positionResponse struct {
		Status   int        `json:"Status"`
		Resource string     `json:"Resource"`
		Payload  []Position `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&positionResponse); err != nil {
		return nil, fmt.Errorf("error decoding position response: %w", err)
	}

	// Return the position information
	return positionResponse.Payload, nil
}
