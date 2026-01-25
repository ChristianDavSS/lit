package domain

type Analyzer interface {
	AnalyzeFile(filePath string, code *[]string) []*FunctionData
}

type ConfigAdapter interface {
	GetConfig() *Config
}

type CacheStore[T comparable, R any] interface {
	GetCache(T) (R, bool)
	SetCache(T, R)
}
