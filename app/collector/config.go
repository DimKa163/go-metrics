package collector

type Config struct {
	Addr           string
	ReportInterval int
	PollInterval   int
	Key            string
	Limit          int
}
