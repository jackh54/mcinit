package server

// RCON client for graceful server shutdown
// This is a placeholder for future implementation

type RCONClient struct {
	host     string
	port     int
	password string
}

func NewRCONClient(host string, port int, password string) *RCONClient {
	return &RCONClient{
		host:     host,
		port:     port,
		password: password,
	}
}

func (r *RCONClient) SendCommand(command string) (string, error) {
	// TODO: Implement RCON protocol
	// For now, this is just a placeholder
	return "", nil
}

func (r *RCONClient) Stop() error {
	_, err := r.SendCommand("stop")
	return err
}

