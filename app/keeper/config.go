package keeper

type Config struct {
	Addr          string
	Path          string
	StoreInterval int64
	Restore       bool
	LogLevel      string
	DatabaseDSN   string
	Key           string
}
