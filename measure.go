package kis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Measure represents the Measure information returned by the Kubota API.
type Measure struct {
	MachineUUID  string     `json:"MachineUUID"`
	MeasureName  string     `json:"MeasureName"`
	MeasureUnit  string     `json:"MeasureUnit"`
	MeasureValue float64    `json:"MeasureValue"`
	Timestamp    CustomTime `json:"Timestamp"`
	CreateTime   CustomTime `json:"CreateTime"`
}

// GetHistoricalMeasureByMobilePhone retrieves historical measure information by mobile phone number.
func (k *Kubota) GetHistoricalMeasureByMobilePhone(mobilePhone, subscription string, startDate, endDate time.Time) ([]Measure, error) {
	return k.getMeasure("mobilePhone", mobilePhone, subscription, startDate, endDate)
}

// GetHistoricalMeasureByUserName retrieves historical measure information by username.
func (k *Kubota) GetHistoricalMeasureByUserName(userName, subscription string, startDate, endDate time.Time) ([]Measure, error) {
	return k.getMeasure("userName", userName, subscription, startDate, endDate)
}

// GetHistoricalMeasureByMachineUUID retrieves historical measure information by machine UUID.
func (k *Kubota) GetHistoricalMeasureByMachineUUID(machineUUID, subscription string, startDate, endDate time.Time) ([]Measure, error) {
	return k.getMeasure("machineUUID", machineUUID, subscription, startDate, endDate)
}

// getMeasure is a helper function to retrieve measure information based on a given field.
func (k *Kubota) getMeasure(field, value, subscription string, startDate, endDate time.Time) ([]Measure, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/measure?%s=%s", k.authentication.Endpoint, field, value)
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
		return nil, fmt.Errorf("error creating measure request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making measure request: %w", err)
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
	var measureResponse struct {
		Status   int       `json:"Status"`
		Resource string    `json:"Resource"`
		Payload  []Measure `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&measureResponse); err != nil {
		return nil, fmt.Errorf("error decoding measure response: %w", err)
	}

	// Return the measure information
	return measureResponse.Payload, nil
}
