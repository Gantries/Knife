package tel

type Protocol string

const (
	GRPC Protocol = "grpc"
	HTTP Protocol = "http"
)

type Remote struct {
	Protocol Protocol `yaml:"proto" default:"http"`
	Endpoint string   `yaml:"endpoint"`
}

type OpenTelemetry struct {
	Metric Remote `yaml:"metric"`
	Trace  Remote `yaml:"trace"`
	Log    Remote `yaml:"log"`
}
