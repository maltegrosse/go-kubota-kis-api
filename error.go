package kis

// errorResponse represents the errorResponse information returned by the Kubota API.
type errorResponse struct {
	Type    string   `json:"Type"`
	Title   string   `json:"Title"`
	Status  int      `json:"Status"`
	LogID   string   `json:"LogId"`
	Details []string `json:"Details"`
}
