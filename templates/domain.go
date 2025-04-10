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
	ports_out "{{.DomainRoot}}/ports/out"
	ports_in "{{.DomainRoot}}/ports/in"
)

type {{.UsecaseName}}Usecase struct {
	Repo ports_out.ModelRepository

	// Add other ports here
}

func (u *{{.UsecaseName}}Usecase) Execute(cmd ports_in.{{.UsecaseName}}Command) error {
	// Add usecase logic here
	return nil
}

`

const DomainExampleInPort = //
`package {{.DomainName}}_ports_in

type {{.UsecaseName}}Command struct {
	// Add command fields here
}

type {{.UsecaseName}}InputPort interface {
	Execute(cmd {{.UsecaseName}}Command) error
}

`

const DomainExampleOutPort = //
`package {{.DomainName}}_ports_out

// Example out port of your domain

import model "{{.DomainRoot}}/model"

type ModelRepository interface {
	Save(model *model.Model) error
}

`

const DomainExampleInAdapter = //
`package {{.DomainName}}_adapter_in

// Example in adapter of your domain

import (
	ports_in "{{.DomainRoot}}/ports/in"
	"net/http"
)

type HttpHandler struct {
	Usecase ports_in.ExampleInputPort
}

func (hh *HttpHandler) Get(w http.ResponseWriter, r *http.Request) {
	hh.Usecase.Execute(ports_in.ExampleCommand{})
}

`

const DomainExampleOutAdapter = //
`package {{.DomainName}}_adapter_out

// Example out adapter of your domain

import model "{{.DomainRoot}}/model"

type PostgresRepository struct {
	// connection ...
}

func (pr *PostgresRepository) Save(model *model.Model) error {
	// Insert into database ...
	return nil
}

`

const DomainExampleWiring = //
`package {{.DomainName}}

// Example wiring of your domain

import (
	adapter_in "{{.DomainRoot}}/adapter/in"
	adapter_out "{{.DomainRoot}}/adapter/out"
	usecase "{{.DomainRoot}}/usecase"
)

func Wiring() {
	repo := adapter_out.PostgresRepository{}
	usecase := usecase.GetModelUsecase{
		Repo: repo,
	}
	handler := adapter_in.HttpHandler{
		Usecase: usecase,
	}

	// handler ...
}

`
