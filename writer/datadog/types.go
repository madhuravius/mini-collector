package datadog

type Config struct {
	ApiKey string `json:"api_key"`
}

type datadogPoint = []interface{}

type datadogSeries struct {
	Metric string         `json:"metric"`
	Points []datadogPoint `json:"points"`
	Type   string         `json:"type"`
	Tags   []string       `json:"tags"`
}

type datadogPayload struct {
	Series []datadogSeries `json:"series"`
}
