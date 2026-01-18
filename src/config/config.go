package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Config struct used to code/decode the JSON file.
type Config struct {
	NamingConvention string `json:"activeNamingConvention"`
}

var ActiveNamingConvention string
var ConfigPath string

// osClear - > Takes the name of the os and returns the command lines to clear the stdout.
var osClear = map[string][]string{
	"windows": []string{"cmd", "/c", "cls"},
	"linux":   []string{"clear"},
}

// Init - > Entry point: Configure all the rules
func Init() {
	newConfig := Config{
		NamingConvention: getNamingConvention(),
	}
	SetNewConfig(newConfig)
	fmt.Println("Configuration updated.")
}

// Function used to clear the screen in Windows or Linux...
func clearScreen(name string, args ...string) {
	// Clear the screen
	cmd := exec.Command(name, args...)
	// Set the stdout we want to clear
	cmd.Stdout = os.Stdout
	// Run the command
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

// SetNewConfig - > Sets up a new config in the json file.
func SetNewConfig(newConfig Config) {
	byteData, err := json.Marshal(newConfig)
	if err != nil {
		os.Exit(1)
	}
	err = os.WriteFile(ConfigPath, byteData, 0644)
	if err != nil {
		os.Exit(1)
	}
}

// LoadDefaultConfig - > Used to load the default values from the JSON config file.
func LoadDefaultConfig() {
	file := ReadConfigJSON()
	var res Config
	err := json.Unmarshal(file, &res)
	if err != nil {
		os.Exit(1)
	}
	// Set the file config as the current config
	ActiveNamingConvention = res.NamingConvention
}

func ReadConfigJSON() []byte {
	path, err := os.Executable()
	if err != nil {
		os.Exit(1)
	}

	ConfigPath = path[:len(path)-7] + "config.json"
	file, err := os.ReadFile(ConfigPath)
	if err != nil {
		os.Exit(1)
	}
	return file
}
