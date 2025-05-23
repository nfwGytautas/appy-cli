package variant_hmd

const templateMainGo = //
`package main

// Auto-generated by appy don't modify by hand, changes will be lost

import (
	"os"
	"os/signal"
	"syscall"

	"{{.Config.Module}}/domains"
	"{{.Config.Module}}/providers"
)

func main() {
	p, err := providers.Initialize()
	if err != nil {
		panic(err)
	}

	err = domains.RegisterDomains(p)
	if err != nil {
		panic(err)
	}

	err = providers.Start(p)
	if err != nil {
		panic(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
`

const templateProvidersGo = //
`package providers

// Auto-generated by appy don't modify by hand, changes will be lost

import (
	{{range .Providers}}providers_{{.Name}} "{{$.Module}}/providers/{{.Name}}"
	{{end}}
)

type Providers struct {
    {{range .Providers}}{{TitleString .Name}} *providers_{{.Name}}.Provider
    {{end}}
}

func Initialize() (*Providers, error) {
	providers := &Providers{}

	{{if .Providers}}
	var err error
	{{end}}

	{{range .Providers}}providers.{{TitleString .Name}}, err = providers_{{.Name}}.Init()
    if err != nil {
        return nil, err
    }
    {{end}}

	return providers, nil
}

func Start(providers *Providers) error {
	{{if .Providers}}
	var err error
	{{end}}

	{{range .Providers}}err = providers.{{TitleString .Name}}.Start()
    if err != nil {
        return err
    }
    {{end}}

	return nil
}
`

const templateDomainsGo = //
`package domains

import (
	"{{.Config.Module}}/domains/example"
)

func RegisterDomains(providers *providers.Providers) error {
	// Adapters
	// ...

	// Domains
	example := example.ExampleDomain{}

	// Connectors
	// ...

	return nil
}
`

const templateGoMod = //
`module {{.Config.Module}}

go 1.23

require github.com/nfwGytautas/appy-go v0.2.0
`

const templateReadmeMd = //
`# {{.Config.Project}}

## Changelog

### 0.1.0 [xxxx/xx/xxx]

#### Changes

- TODO

#### Bugfixes

- N/A
`

const templateDockerfile = //
`FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

FROM alpine
COPY --from=build /server /server
EXPOSE 3000
ENTRYPOINT ["/server"]
`

const templateGitignore = //
`# Go
bin/
build/
dist/

# Appy
.appy/

# General
.env
.DS_Store
`

const templateVscodeSnippets = //
`{
    "Changelog": {
        "prefix": "changelog",
        "scope": "Markdown",
        "body": [
            "### ${1:0.0.0} [${CURRENT_YEAR}/${CURRENT_MONTH}/${CURRENT_DATE}]",
            "",
            "#### Changes",
            "",
            "- ${2:Feature}",
            "",
            "#### Bugfixes",
            "",
            "- ${3:N/A}"
        ],
        "description": "Create a changelog entry"
    }
}
`

const templateVscodeSettings = //
`{
    "cSpell.words": [
        "appy",
        "Usecase"
    ]
}
`

const templateGithubBuildYaml = //
`name: Build binary

on:
  pull_request:
    branches:
        - main

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
          with:
            submodules: true
            token: {{"${{ secrets.GITHUB_TOKEN }}"}}

      - uses: actions/setup-go@v4
        with:
            go-version: 1.23

      - run: go build -o /{{.Project}} .

      - uses: actions/upload-artifact@v4
        with:
            name: {{.Config.Project}}
            path: /{{.Config.Project}}

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          submodules: true
          token: {{"${{ secrets.GITHUB_TOKEN }}"}}

      - uses: actions/setup-go@v4
        with:
          go-version: 1.23

      - run: go test ./...
`

const templateErrorsGo = //
`package shared

// Define your program errors here

import "errors"

type Error struct {
    Code string
}

// e.g.
var (
    ErrNotFound = errors.New("not found")
)
`

const templateDomainExampleDomain = //
`package {{.DomainName}}

// Describe the domain in this file add dependencies that will need adapters, etc.

type {{TitleString .DomainName}}Domain struct {
}
`

const templateDomainExampleModel = //
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

const templateDomainExampleUsecase = //
`package {{.DomainName}}

type {{TitleString .UsecaseName}}Args struct {
	// Add usecase arguments here
}

func (d *{{TitleString .DomainName}}Domain) {{TitleString .UsecaseName}}(args {{TitleString .UsecaseName}}Args) error {
	// Add usecase logic here
	return nil
}
`
