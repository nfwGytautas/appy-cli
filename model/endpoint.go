package model

type Middleware struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Params map[string]string `yaml:"params"`
}

type EndpointGroup struct {
	Name       string   `yaml:"name"`
	Path       string   `yaml:"path"`
	Parent     string   `yaml:"parent"`
	Middleware []string `yaml:"middleware"`
}
