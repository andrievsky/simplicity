package config

import "fmt"

const BackendName = "Simplicity"
const BackendVersion = "0.0.1"
const BackendPort = "8090"

func BackendInfo() string {
	return fmt.Sprintf("%s %s", BackendName, BackendVersion)
}
