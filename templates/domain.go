package templates

const DomainExampleModel = //
`package {{.DomainName}}_model

// Example model of your domain

type Model struct {
	ID string
}

func NewModel(id string) *Model {
	return &Model{
		ID: id,
	}
}

`

const DomainExampleUsecase = //
`package {{.DomainName}}_usecase

// Example usecase of your domain

import (
	ports "{{.DomainRoot}}/ports"
)

type {{.UsecaseName}}Usecase struct {
	Repo ports.ModelRepository

	// Add other ports here
}

func (u *{{.UsecaseName}}Usecase) Execute(cmd ports.{{.UsecaseName}}Command) error {
	// Add usecase logic here
	return nil
}

`

const DomainExampleInPort = //
`package {{.DomainName}}_ports_in

// An input port is something that your domain can do

type {{.UsecaseName}}Command struct {
	// Add command fields here
}

type {{.UsecaseName}}InputPort interface {
	Execute(cmd {{.UsecaseName}}Command) error
}

`

const DomainExampleOutPort = //
`package {{.DomainName}}_ports_out

// An output port is something that your domain needs

import model "{{.DomainRoot}}/model"

type ModelRepository interface {
	Save(model *model.Model) error
}

`
