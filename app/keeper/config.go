package keeper

type Config struct {
	Addr               string `arg:"a" envArg:"ADDRESS" json:"address"`
	Path               string `arg:"f" envArg:"FILE_STORAGE_PATH" json:"store_file"`
	StoreInterval      int64  `arg:"i" envArg:"STORE_INTERVAL" json:"store_interval"`
	Restore            bool   `arg:"r" envArg:"RESTORE" json:"restore"`
	LogLevel           string `arg:"l" envArg:"LOG_LEVEL"`
	DatabaseDSN        string `arg:"d" envArg:"DATABASE_DSN" json:"database_dsn"`
	Key                string `arg:"k" envArg:"KEY" json:"key"`
	PrivateKeyFilePath string `arg:"c" envArg:"CRYPTO_KEY" json:"crypto_key"`
}
