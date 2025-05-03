package config

type Config struct {
	BackendName    string `json:"backend_name"`
	BackendVersion string `json:"backend_version"`
	AWS            AWS    `json:"aws"`
	Server         Server `json:"server"`
	EnableDebug    bool   `json:"debug"`
}

type Server struct {
	Port string `json:"port"`
}

type AWS struct {
	Profile string `json:"profile"`
	Bucket  string `json:"bucket"`
}

func LoadConfig() (*Config, error) {
	config := &Config{
		BackendName:    "Simplicity",
		BackendVersion: "0.4.0",
		Server: Server{
			Port: "8090",
		},
		AWS: AWS{
			Profile: "nick-aws-personal",
			Bucket:  "simplicity-backend-storage",
		},
		EnableDebug: false,
	}
	return config, nil
}
