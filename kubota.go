package kis

// kubota represents the Kubota API client.
type kubota struct {
	authentication *authentication
}

// NewKIS creates a new Kubota API client.
func NewKIS(publicKey, SecretKey, Endpoint string) (*kubota, error) {
	k := &kubota{}
	auth, err := newAuthentication(publicKey, SecretKey, Endpoint)
	if err != nil {
		return nil, err
	}
	k.authentication = auth
	return k, nil
}
