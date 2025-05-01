package server

type Config struct {
	CORS        bool // Enable CORS
	Security    bool // Enable security headers and server middleware
	GatewayMode bool
	AllowOrigin string
}
