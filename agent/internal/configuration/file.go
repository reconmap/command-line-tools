package configuration

type KeycloakConfig struct {
	BaseUri      string `json:"baseUri"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type ReconmapApiConfig struct {
	BaseUri string `json:"baseUri"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	KeycloakConfig    `json:"keycloak"`
	ReconmapApiConfig `json:"reconmapApi"`
	RedisConfig       `json:"redis"`
	ValidOrigins      string `json:"validOrigins"`
}

const ConfigFileName string = "config-reconmapd.json"

func NewConfig() Config {
	return Config{
		KeycloakConfig: KeycloakConfig{
			BaseUri:      "http://localhost:8080",
			ClientID:     "reconmapd-cli",
			ClientSecret: "REPLACE THIS WITH YOUR CLIENT SECRET",
		},
		ReconmapApiConfig: ReconmapApiConfig{
			BaseUri: "http://localhost:5510",
		},
		RedisConfig: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "REconDIS",
		},
		ValidOrigins: "http://localhost:5500",
	}
}
