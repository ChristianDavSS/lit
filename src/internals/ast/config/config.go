package config

import "CLI_App/src/config"

// ast/config is a middleware for variables between config/ and ast/, acceding to them over here instead
// of importing the global config itself

var ActivePattern, ActiveConventionIndex = config.GetActiveNamingConvention()

// ShouldFix flag to indicate if we should fix the naming of the variables
var ShouldFix = false
