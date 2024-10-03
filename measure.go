package kis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
func (k *kubota) GetHistoricalMeasureByMobilePhone(mobilePhone, subscription string, startDate, endDate time.Time) ([]Measure, error) {
	return k.getMeasure("mobilePhone", mobilePhone, subscription, startDate, endDate)
}

// GetHistoricalMeasureByUserName retrieves historical measure information by username.
func (k *kubota) GetHistoricalMeasureByUserName(userName, subscription string, startDate, endDate time.Time) ([]Measure, error) {
	return k.getMeasure("userName", userName, subscription, startDate, endDate)
}

// GetHistoricalMeasureByMachineUUID retrieves historical measure information by machine UUID.
func (k *kubota) GetHistoricalMeasureByMachineUUID(machineUUID, subscription string, startDate, endDate time.Time) ([]Measure, error) {
	return k.getMeasure("machineUUID", machineUUID, subscription, startDate, endDate)
}

// getMeasure is a helper function to retrieve measure information based on a given field.
func (k *kubota) getMeasure(field, value, subscription string, startDate, endDate time.Time) ([]Measure, error) {
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting measure: %s", string(body))
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
