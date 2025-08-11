package services

import(
	"os"
	"gopkg.in/yaml.v3"
)



type Config struct {
	Site struct {
		Title string `yaml:"title"`
	} `yaml:"site"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
}


func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}


	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err 
	}

	
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err 
	}

	return config, nil
}