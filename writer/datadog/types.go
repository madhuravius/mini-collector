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
	Host   string         `json:"host"`
	Tags   []string       `json:"tags"`
}

type datadogPayload struct {
	Series []datadogSeries `json:"series"`
}
