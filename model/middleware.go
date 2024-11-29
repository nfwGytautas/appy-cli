package model

type Middleware struct {
	Name     string            `yaml:"name"`
	Provider string            `yaml:"provider"`
	Params   map[string]string `yaml:"params"`
}
