package model

type EndpointGroup struct {
	Name       string   `yaml:"name"`
	Path       string   `yaml:"path"`
	Parent     string   `yaml:"parent"`
	Middleware []string `yaml:"middleware"`
}

type Endpoint struct {
	Name     string     `yaml:"name"`
	Group    string     `yaml:"group"`
	Method   string     `yaml:"method"`
	Impl     string     `yaml:"impl"`
	Children []Endpoint `yaml:"children"`
}

func (e *Endpoint) ResolveChildren() {
	for i := range e.Children {
		if e.Children[i].Group == "" {
			e.Children[i].Group = e.Group
		}

		e.Children[i].Name = e.Name + "/" + e.Children[i].Name
		e.Children[i].ResolveChildren()
	}
}
