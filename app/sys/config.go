package sys

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Configuration struct {
	Name     string `json:"Name"`
	Settings struct {
		Language       string `json:"Language"`
		TemplateEngine string `json:"TemplateEngine"`
		BaseUrl        string `json:"BaseUrl"`
		AutoRoutes     bool   `json:"AutoRoutes"`
		UseViewOutput  bool   `json:"UseViewOutput"`
		StatsEnabled   bool   `json:"StatsEnabled"`
	}
	Database struct {
		Name      string `json:"Name"`
		Host      string `json:"Host"`
		Port      int    `json:"Port"`
		UserName  string `json:"UserName"`
		Password  string `json:"Password"`
		DBEnabled bool   `json:"DBEnabled"`
	}
}

var Config = GetConfig()

func GetConfig() Configuration {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Read the file content
	byteValue, _ := io.ReadAll(jsonFile)

	// Unmarshal the JSON data
	var config Configuration
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}
	//	fmt.Printf("Type: %+v\n", config.Settings.Preferences.Language)
	//	fmt.Printf("Name: %+v\n", config.Settings.Preferences.TemplateEngine)
	return config
}
