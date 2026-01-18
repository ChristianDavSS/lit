package main

import (
	root "CLI_App/src/cmd"
	"CLI_App/src/config"
)

func main() {
	// Load up the config from the json file.
	config.LoadDefaultConfig()
	root.Execute()
}
