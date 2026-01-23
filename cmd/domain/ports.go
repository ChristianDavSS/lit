package domain

type Analyzer interface {
	AnalyzeFile(filePath string, code []byte) []*FunctionData
}

type ConfigAdapter interface {
	GetConfig() *Config
	SaveConfig(cfg *Config)
}
