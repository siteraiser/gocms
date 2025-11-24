package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Name     string `json:"Name"`
	Settings struct {
		Preferences struct {
			Language       string `json:"Language"`
			TemplateEngine string `json:"TemplateEngine"`
			BaseUrl        string `json:"BaseUrl"`
		}
	}
}

func GetConfig() Config {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Read the file content
	byteValue, _ := io.ReadAll(jsonFile)

	// Unmarshal the JSON data
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}

	//	fmt.Printf("Type: %+v\n", config.Settings.Preferences.Language)
	//	fmt.Printf("Name: %+v\n", config.Settings.Preferences.TemplateEngine)
	return config
}
