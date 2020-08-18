package conf

type Config struct {
	Core      Core      `json:"core" yaml:"core"`
	Level     Level     `json:"level" yaml:"level"`
	Formatter Formatter `json:"formatter" yaml:"formatter"`
	Outputs   []Output  `json:"outputs" yaml:"outputs"`
}
