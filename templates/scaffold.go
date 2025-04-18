package templates

const AppyYaml = //
`version: {{.Version}}
project: {{.ProjectName}}
type: {{.Type}}
module: {{.Module}}
repositories:
    - url: https://github.com/nfwGytautas/appy-providers
      branch: main
`

const MainGo = //
`package main

// Auto-generated by appy don't modify by hand, changes will be lost

type Providers struct {}

func main() {
	var err error
	providers := &Providers{}

	domains, err := connectDomains(providers)
	if err != nil {
		panic(err)
	}

	start(providers, domains)
}
`

const MainWithProvidersGo = //
`package main

// Auto-generated by appy don't modify by hand, changes will be lost

import (
	{{range .Providers}}providers_{{.Name}} "{{$.Module}}/providers/{{.Name}}"{{end}}
)

type Providers struct {
	{{range .Providers}}{{.Name}} *providers_{{.Name}}.Provider{{end}}
}

func main() {
	var err error
	providers := &Providers{}

	{{range .Providers}}providers.{{.Name}}, err = providers_{{.Name}}.Init()
	if err != nil {
		panic(err)
	}{{end}}

	domains, err := connectDomains(providers)
	if err != nil {
		panic(err)
	}

	start(providers, domains)
}
`

const WiringGo = //
`package main

type Domains struct {}

func connectDomains(providers *Providers) (*Domains, error) {
	// Connect your domains here

	return &Domains{}, nil
}

func start(providers *Providers, domains *Domains) {
	// If any provider has a "Start" method, call it here
}

`

const GoMod = //
`module {{.Module}}

go 1.23

require github.com/nfwGytautas/appy-go v0.1.0
`

const ReadmeMd = //
`# {{.ProjectName}}

## Changelog

### 0.1.0 [xxxx/xx/xxx]

#### Changes

- TODO

#### Bugfixes

- N/A
`

const Dockerfile = //
`FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN go build -o /server .

FROM scratch
COPY --from=build /server /server
EXPOSE 3000
CMD ["/server"]
`

const Gitignore = //
`# Go
bin/
build/
dist/
`

const VscodeSnippets = //
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

const VscodeSettings = //
`{
    "cSpell.words": [
        "appy"
    ]
}
`

const GithubBuildYaml = //
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

	  - run: go build -o /{{.ProjectName}} .

	  - uses: actions/upload-artifact@v4
		with:
			name: {{.ProjectName}}
			path: /{{.ProjectName}}
`

const ErrorsGo = //
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
