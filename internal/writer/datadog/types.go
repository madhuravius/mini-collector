package datadog

type Config struct {
	ApiKey     string `json:"api_key"`
	Timeout    string `json:"timeout"`
	RetryCount string `json:"retry_count"`
	SeriesUrl  string `json:"series_url"`
}

type datadogPoint = []interface{}

type datadogSeries struct {
	Metric string         `json:"metric"`
	Points []datadogPoint `json:"points"`
	Type   string         `json:"type"`
	// Host directly collides with the http spec, due to the need
	// for GRPC to submit both :authority and host, we will lose
	// this header. In some places we will use .host_name, which
	// this is also an alias for.
	Host string   `json:"host"`
	Tags []string `json:"tags"`
}

type datadogPayload struct {
	Series []datadogSeries `json:"series"`
}
