package collector

type Config struct {
	Addr              string `arg:"a" envArg:"ADDRESS" json:"address"`
	ReportInterval    int    `arg:"r" envArg:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval      int    `arg:"p" envArg:"POLL_INTERVAL" json:"poll_interval"`
	Key               string `arg:"k" envArg:"KEY" json:"key"`
	Limit             int    `arg:"r" envArg:"RATE_LIMIT" json:"rate_limit"`
	PublicKeyFilePath string `arg:"c" envArg:"CRYPTO_KEY" json:"crypto_key"`
}
