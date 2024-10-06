package kis

// Error represents the Error information returned by the Kubota API.
type Error struct {
	Type    string   `json:"Type"`
	Title   string   `json:"Title"`
	Status  int      `json:"Status"`
	LogID   string   `json:"LogId"`
	Details []string `json:"Details"`
}