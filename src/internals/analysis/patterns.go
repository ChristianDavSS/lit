package analysis

// patterns.go - > File that contains every regex pattern used in the script

// String scanValidScriptPattern for the regex that validates scripts for the scanner.
// Only contains the languages supported for scanning
var scanValidScriptPattern = "^[a-zA-Z0-9._-]+\\.(py|js|java)$"

// String locValidStringPattern validates the scripts for the loc.
var locValidScriptPattern = "^[a-zA-Z0-9._-]+\\.(py|go|java|js|jsx|dart|c|cpp|css|html|ts|md)$"

// String notValidDirPattern for the regex that validates you won´t visit unwanted sites.
// Only contains unwanted directories and files.
var notValidDirPattern = "^(node_modules|.*\\.exe|target|venv|__pycache__|" +
	"\\.(git|idea|mvn|cmd))$"
