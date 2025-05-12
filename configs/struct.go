package configs

type AppConfig struct {
	FilePath    string
	Strategy    string
	ChunkSize   int
	WorkerCount int
	LogLevel    string
	EnableProf  bool
}

const (
	StrategyConcurrent = "concurrent"
	StrategyLinear     = "linear"
)
