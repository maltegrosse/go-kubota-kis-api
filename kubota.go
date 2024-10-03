package kis

// Kubota represents the Kubota API client.
type Kubota struct {
	authentication *authentication
}

// NewKIS creates a new Kubota API client.
func NewKIS(publicKey, SecretKey, Endpoint string) (*Kubota, error) {
	k := &Kubota{}
	auth, err := newAuthentication(publicKey, SecretKey, Endpoint)
	if err != nil {
		return nil, err
	}
	k.authentication = auth
	return k, nil
}
