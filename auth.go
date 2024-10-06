package kis

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type authentication struct {
	PublicKey   string
	SecretKey   string
	Endpoint    string
	Token       string
	TokenType   string
	ExpiresIn   int
	TokenExpiry time.Time
	// Mutex to protect the token during refresh
	TokenMutex sync.Mutex
}

// tokenResponse represents the JSON response from the token endpoint.
type tokenResponse struct {
	AccessToken string `json:"AccessToken"`
	TokenType   string `json:"TokenType"`
	ExpiresIn   int    `json:"ExpiresIn"`
}

func newAuthentication(publicKey, SecretKey, Endpoint string) (*authentication, error) {
	a := &authentication{
		PublicKey: publicKey,
		SecretKey: SecretKey,
		Endpoint:  Endpoint,
	}
	err := a.accessToken()
	if err != nil {
		return nil, errors.Join(err, errors.New("error with authentication"))
	}
	go a.refreshAccessToken()

	return a, nil
}
func (a *authentication) getToken() string {
	a.TokenMutex.Lock()
	t := a.Token
	defer a.TokenMutex.Unlock()
	return t
}

// refreshToken retrieves a new access token from the Kubota API.
func (a *authentication) accessToken() error {
	a.TokenMutex.Lock()
	tokenURL := a.Endpoint
	// Construct the request URL and payload
	data := url.Values{}
	data.Set("grantType", "authorization_code")
	data.Set("publicKey", a.PublicKey)
	data.Set("secretKey", a.SecretKey)

	u, err := url.Parse(tokenURL)
	if err != nil {
		return fmt.Errorf("error parsing token URL: %w", err)
	}
	u.Path = "/api/v1/authorization/token"
	// Make the request
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making token request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		var errResponse = errorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return fmt.Errorf("error decoding error response: %w", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			// get the header and retry after the specified time
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				errResponse.Details = append(errResponse.Details, fmt.Sprintf("Retry-After: %s", retryAfter))
			}
		}
		return fmt.Errorf("error: %s with statuscode: %d, type %s, details: %s", errResponse.Title, errResponse.Status, errResponse.Type, strings.Join(errResponse.Details, ", "))
	}

	// Unmarshal the response
	var tokenResponse tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("error decoding token response: %w", err)
	}

	// Update the configuration with the new token information
	a.Token = tokenResponse.AccessToken
	a.TokenType = tokenResponse.TokenType
	a.ExpiresIn = tokenResponse.ExpiresIn
	a.TokenExpiry = time.Now().Add(time.Duration(a.ExpiresIn) * time.Minute)
	a.TokenMutex.Unlock()
	return nil
}

// refreshAccessToken periodically refreshes the access token.
func (a *authentication) refreshAccessToken() {
	for {
		// check if token is 1 minute before expiry
		if time.Now().After(a.TokenExpiry.Add(-1 * time.Minute)) {

			// Refresh the token
			if err := a.accessToken(); err != nil {
				// Handle refresh errors
				log.Fatal(errors.Join(err, errors.New("error refreshing access token")))
			} else {
				// Log the successful refresh
				log.Println("Access token refreshed")
			}
		}
		// Sleep until the next refresh attempt, using the ExpiresIn value
		time.Sleep(time.Duration(a.ExpiresIn-1) * time.Minute) // Subtract 1 minute for a safety margin
	}
}
