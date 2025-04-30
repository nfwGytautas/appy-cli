package templates

const DomainExampleDomain = //
`package {{.DomainName}}

// Describe the domain in this file add dependencies that will need adapters, etc.

type {{TitleString .DomainName}}Domain struct {
}
`

const DomainExampleModel = //
`package {{.DomainName}}_model

type {{TitleString .DomainName}} struct {
	ID string
}

func New{{TitleString .DomainName}}(id string) *{{TitleString .DomainName}} {
	return &{{TitleString .DomainName}}{
		ID: id,
	}
}
`

const DomainExampleUsecase = //
`package {{.DomainName}}

type {{TitleString .UsecaseName}}Args struct {
	// Add usecase arguments here
}

func (d *{{TitleString .DomainName}}Domain) {{.UsecaseName}}(args {{TitleString .UsecaseName}}Args) error {
	// Add usecase logic here
	return nil
}
`
