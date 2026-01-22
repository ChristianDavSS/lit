package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Config struct used to code/decode the JSON file.
type Config struct {
	NamingConventionIndex int8 `json:"activeNamingConventionIndex"`
}

var ActiveConfig *Config
var ConfigPath string

// osClear - > Takes the name of the os and returns the command lines to clear the stdout.
var osClear = map[string][]string{
	"windows": []string{"cmd", "/c", "cls"},
	"linux":   []string{"clear"},
}

// Init - > Entry point: Configure all the rules
func Init() {
	newConfig := Config{
		NamingConventionIndex: getNamingConvention(),
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
		fmt.Fprintln(os.Stderr, "Error running the command to clear your screen.")
		os.Exit(1)
	}
}

// SetNewConfig - > Sets up a new config in the json file.
func SetNewConfig(newConfig Config) {
	byteData, err := json.Marshal(newConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing the struct data into json.")
		os.Exit(1)
	}
	err = os.WriteFile(ConfigPath, byteData, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing the data into the config file...")
		os.Exit(1)
	}
}

// LoadDefaultConfig - > Used to load the default values from the JSON config file.
func LoadDefaultConfig() {
	file := ReadConfigJSON()
	err := json.Unmarshal(file, &ActiveConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading the config file...")
		os.Exit(1)
	}
}

func ReadConfigJSON() []byte {
	path, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting the .exe path name")
		os.Exit(1)
	}

	ConfigPath = filepath.Dir(path) + "/config.json"
	file, err := os.ReadFile(ConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't read the config file (is it in '"+ConfigPath+"'?)")
		os.Exit(1)
	}
	return file
}

// GetActiveNamingConvention function: gets the active naming convention from the current configuration file
func GetActiveNamingConvention() (string, int8) {
	// If the 'ActiveConfig' variable is nil, we set a value to it
	if ActiveConfig == nil {
		LoadDefaultConfig()
	}

	return conventions[ActiveConfig.NamingConventionIndex], ActiveConfig.NamingConventionIndex
}
