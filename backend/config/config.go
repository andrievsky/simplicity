package config

const BackendName = "Simplicity"
const BackendVersion = "0.3.1"
const BackendPort = "8090"

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func BackendInfo() Info {
	return Info{
		Name:    BackendName,
		Version: BackendVersion,
	}
}
