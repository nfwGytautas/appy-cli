package config

type Repository struct {
	Url    string `yaml:"url"`
	Branch string `yaml:"branch"`
}
