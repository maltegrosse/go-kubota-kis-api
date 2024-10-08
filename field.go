package kis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Field represents the Field information returned by the Kubota API.
type Field struct {
	FieldID     string     `json:"FieldID"`
	CompanyID   string     `json:"CompanyID"`
	FieldName   string     `json:"FieldName"`
	Shape       Shape      `json:"Shape"`
	FieldStatus string     `json:"FieldStatus"`
	Timestamp   CustomTime `json:"Timestamp"`
	CreateTime  CustomTime `json:"CreateTime"`
	UpdateTime  CustomTime `json:"UpdateTime"`
}

// Shape represents the GeoJSON Shape of the field.
type Shape struct {
	Type        string         `json:"type"`
	Coordinates [][][2]float64 `json:"coordinates"`
}

// GetFieldByMobilePhone retrieves field information by mobile phone number.
func (k *Kubota) GetFieldByMobilePhone(mobilePhone string) ([]Field, error) {
	return k.getField("mobilePhone", mobilePhone)
}

// GetFieldByUserName retrieves field information by username.
func (k *Kubota) GetFieldByUserName(userName string) ([]Field, error) {
	return k.getField("userName", userName)
}

// getField is a helper function to retrieve field information based on a given field.
func (k *Kubota) getField(f, value string) ([]Field, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/field?%s=%s", k.authentication.Endpoint, f, value)

	// Make the request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating field request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making field request: %w", err)
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
	var fieldResponse struct {
		Status   int     `json:"Status"`
		Resource string  `json:"Resource"`
		Payload  []Field `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&fieldResponse); err != nil {
		return nil, fmt.Errorf("error decoding field response: %w", err)
	}

	// Return the field information
	return fieldResponse.Payload, nil
}
