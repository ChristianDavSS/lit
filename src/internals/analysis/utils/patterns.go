package utils

// patterns.go - > File that contains every regex pattern used in the script

// ScanValidScriptPattern for the regex that validates scripts for the scanner.
// Only contains the languages supported for scanning
var ScanValidScriptPattern = "^[a-zA-Z0-9._-]+\\.(py|js|java|jsx|go)$"

// LocValidScriptPattern validates the scripts for the loc.
var LocValidScriptPattern = "^[a-zA-Z0-9._-]+\\.(py|go|java|js|jsx|dart|c|cpp|css|html|ts|md)$"

// NotValidDirPattern for the regex that validates you won´t visit unwanted sites.
// Only contains unwanted directories and files.
var NotValidDirPattern = "^(node_modules|.*\\.exe|target|venv|__pycache__|" +
	"\\.(git|idea|mvn|cmd))$"

// RemoteUserEmail pattern to detect emails of commits done in remote
var RemoteUserEmail = ".+@users.noreply.github.com"
