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

type GetModelUsecase struct {
	Repo {{.DomainName}}_ports_out.ModelRepository

	// Add other ports here
}

func (gmu *GetModelUsecase) Execute(cmd GetModelCommand) (*{{.DomainName}}_model.Model, error) {
	return gmu.Repo.Get(cmd)
}

`

const DomainExampleInPort = //
`package {{.DomainName}}_ports_in

// Example in port of your domain

type GetModelCommand struct {
	ID string
}

type GetModelInputPort interface {
	Execute(cmd GetModelCommand) (*{{.DomainName}}_model.Model, error)
}

`

const DomainExampleOutPort = //
`package {{.DomainName}}_ports_out

// Example out port of your domain

type ModelRepository interface {
	Save(model *{{.DomainName}}_model.Model) error
}

`

const DomainExampleInAdapter = //
`package {{.DomainName}}_adapter_in

// Example in adapter of your domain

type HttpHandler struct {
	Usecase {{.DomainName}}_ports_in.GetModelInputPort
}

func (hh *HttpHandler) Get(w http.ResponseWriter, r *http.Request) {
	hh.Usecase.Execute({{.DomainName}}_usecase.GetModelCommand{ID: r.URL.Query().Get("id")})
}

`

const DomainExampleOutAdapter = //
`package {{.DomainName}}_adapter_out

// Example out adapter of your domain

type PostgresRepository struct {
	// connection ...
}

func (pr *PostgresRepository) Save(model *{{.DomainName}}_model.Model) error {
	// Insert into database ...
	return nil
}

`

const DomainExampleWiring = //
`package {{.DomainName}}

// Example wiring of your domain

func Wiring() {
	repo := {{.DomainName}}_adapter_out.PostgresRepository{}
	usecase := {{.DomainName}}_usecase.GetModelUsecase{
		Repo: repo,
	}
	handler := {{.DomainName}}_adapter_in.HttpHandler{
		Usecase: usecase,
	}

	// handler ...
}

`
