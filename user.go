package kis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// User represents the User information returned by the Kubota API.
type User struct {
	UserID      string     `json:"UserID"`
	MobilePhone string     `json:"MobilePhone"`
	UserName    string     `json:"UserName"`
	CompanyID   string     `json:"CompanyID"`
	Email       string     `json:"Email"`
	FirstName   string     `json:"FirstName"`
	LastName    string     `json:"LastName"`
	UserStatus  string     `json:"UserStatus"`
	Timestamp   CustomTime `json:"Timestamp"`
	CreateTime  CustomTime `json:"CreateTime"`
	UpdateTime  CustomTime `json:"UpdateTime"`
}

// GetUserByMobilePhone retrieves user information by mobile phone number.
func (k *Kubota) GetUserByMobilePhone(mobilePhone string) (*User, error) {
	return k.getUser("mobilePhone", mobilePhone)
}

// GetUserByUserName retrieves user information by username.
func (k *Kubota) GetUserByUserName(userName string) (*User, error) {
	return k.getUser("userName", userName)
}

// getUser is a helper function to retrieve user information based on a given field.
func (k *Kubota) getUser(field, value string) (*User, error) {
	// Construct the request URL
	apiURL := fmt.Sprintf("%s/api/v1/user?%s=%s", k.authentication.Endpoint, field, value)

	// Make the request
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating user request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.authentication.getToken())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making user request: %w", err)
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
	var userResponse struct {
		Status   int    `json:"Status"`
		Resource string `json:"Resource"`
		Payload  User   `json:"Payload"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error decoding user response: %w", err)
	}

	// Return the user information
	return &userResponse.Payload, nil
}
