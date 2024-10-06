package kis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Alarm represents the Alarm information returned by the Kubota API.
type Alarm struct {
	MachineUUID string     `json:"MachineUUID"`
	Type        string     `json:"Type"`
	Description string     `json:"Description"`
	Timestamp   CustomTime `json:"Timestamp"`
	CreateTime  CustomTime `json:"CreateTime"`
}

// GetHistoricalAlarmByMobilePhone retrieves historical alarm information by mobile phone number.
func (k *Kubota) GetHistoricalAlarmByMobilePhone(mobilePhone, subscription string, startDate, endDate time.Time) ([]Alarm, error) {
	return k.getAlarm("mobilePhone", mobilePhone, subscription, startDate, endDate)
}

// GetHistoricalAlarmByUserName retrieves historical alarm information by username.
func (k *Kubota) GetHistoricalAlarmByUserName(userName, subscription string, startDate, endDate time.Time) ([]Alarm, error) {
	return k.getAlarm("userName", userName, subscription, startDate, endDate)
}

// GetHistoricalAlarmByMachineUUID retrieves historical alarm information by machine UUID.
func (k *Kubota) GetHistoricalAlarmByMachineUUID(machineUUID, subscription string, startDate, endDate time.Time) ([]Alarm, error) {
	return k.getAlarm("machineUUID", machineUUID, subscription, startDate, endDate)
}

// getAlarm is a helper function to retrieve alarm information based on a given field.
func (k *Kubota) getAlarm(field, value, subscription string, startDate, endDate time.Time) ([]Alarm, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/alarm?%s=%s", k.authentication.Endpoint, field, value)
	if subscription != "" {
		apiURL += "&subscription=" + subscription
	}
	if !startDate.IsZero() {
		s := CustomTime{startDate}
		ss, err := s.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshaling start date: %w", err)
		}
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
		return nil, fmt.Errorf("error creating alarm request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making alarm request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		var errResponse = errorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return nil, fmt.Errorf("error decoding error response: %w", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			// get the header and retry after the specified time
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				errResponse.Details = append(errResponse.Details, fmt.Sprintf("Retry-After: %s", retryAfter))
			}
		}
		return nil, fmt.Errorf("error: %s with statuscode: %d, type %s, details: %s", errResponse.Title, errResponse.Status, errResponse.Type, strings.Join(errResponse.Details, ", "))
	}

	// Unmarshal the response
	var alarmResponse struct {
		Status   int     `json:"Status"`
		Resource string  `json:"Resource"`
		Payload  []Alarm `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&alarmResponse); err != nil {
		return nil, fmt.Errorf("error decoding alarm response: %w", err)
	}

	// Return the alarm information
	return alarmResponse.Payload, nil
}
