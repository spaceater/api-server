package config

type ServerConfig struct {
	Port string `json:"port"`
}

var (
	Port string
)

func InitServerConfig(configData map[string]interface{}) {
	Port = getConfigString(getJSONTag(ServerConfig{}, "Port"), configData, "80")
}
