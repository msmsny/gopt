package gopt

func uberGoStyleOptionTemplate() []byte {
	return []byte(`
{{- if .PackageName -}}
package {{.PackageName}}

{{end -}}
{{- if (hasDuration .Options) -}}
import (
	"time"
)

{{end -}}
{{- if .Evaluate -}}
func evaluateOptions(options []{{title .Name}}Option) *{{.Name}}Options {
	opts := &{{.Name}}Options{}
	for _, option := range options {
		option.apply(opts)
	}

	return opts
}

{{end -}}
type {{title .Name}}Option interface{ apply(*{{.Name}}Options) }

type {{.Name}}Options struct {
{{- range $option := .Options}}
	{{$option.Name}} {{$option.Type}}
{{- end}}
}
{{range $option := .Options}}
type {{$option.Name}}Option {{$option.Type}}

func (o {{$option.Name}}Option) apply(opts *{{$.Name}}Options) { opts.{{$option.Name}} = {{$option.Type}}(o) }

func With{{title $option.Name}}(value {{$option.Type}}) {{title $.Name}}Option { return {{$option.Name}}Option(value) }
{{end -}}
`)
}
