package gopt

func uberGoStyleOptionTemplate() []byte {
	return []byte(`
{{- if .Header -}}
// Code generated by gopt. DO NOT EDIT.

{{end -}}
{{- if .PackageName -}}
package {{.PackageName}}

{{end -}}
{{- if (hasImport .Options) -}}
import (
{{- range $option := .Options}}
{{- if ($option.IsDuration) }}
	"time"

{{- else if ($option.IsPackage) -}}
{{- if ($option.Package.IsFullPath) }}
	"{{$option.Package.Base}}/{{$option.Package.Name}}"

{{- else if ($option.Package.HasName) }}
	"{{$option.Package.Name}}"

{{end -}}
{{end -}}
{{end}}
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
{{- if not $option.IsPackage}}
	{{$option.Name}} {{$option.Type}}
{{- else if or $option.Package.IsFullPath $option.Package.HasName}}
	{{$option.Name}} {{$option.Package.Prefix}}{{$option.Package.Name}}.{{$option.Package.Type}}
{{- else}}
	{{$option.Name}} {{$option.Package.Prefix}}{{$option.Package.Type}}
{{- end}}
{{- end}}
}
{{range $option := .Options}}
{{- if not $option.IsPackage}}
type {{$option.Name}}Option {{$option.Type}}

func (o {{$option.Name}}Option) apply(opts *{{$.Name}}Options) { opts.{{$option.Name}} = {{$option.Type}}(o) }

func With{{title $option.Name}}(value {{$option.Type}}) {{title $.Name}}Option { return {{$option.Name}}Option(value) }
{{- else}}
type {{toLower $option.Package.Type}}Option struct {
{{- if or $option.Package.IsFullPath $option.Package.HasName}}
	{{toLower $option.Package.Type}} {{$option.Package.Prefix}}{{$option.Package.Name}}.{{$option.Package.Type}}
{{- else}}
	{{toLower $option.Package.Type}} {{$option.Package.Prefix}}{{$option.Package.Type}}
{{- end}}
}

func (o {{toLower $option.Package.Type}}Option) apply(opts *{{$.Name}}Options) { opts.{{$option.Name}} = o.{{toLower $option.Package.Type}} }
{{- if or $option.Package.IsFullPath $option.Package.HasName}}

func With{{title $option.Name}}({{toLower $option.Package.Type}} {{$option.Package.Prefix}}{{$option.Package.Name}}.{{$option.Package.Type}}) {{title $.Name}}Option { return {{toLower $option.Package.Type}}Option{ {{toLower $option.Package.Type}}: {{toLower $option.Package.Type}} } }
{{- else}}

func With{{title $option.Name}}({{toLower $option.Package.Type}} {{$option.Package.Prefix}}{{$option.Package.Type}}) {{title $.Name}}Option { return {{toLower $option.Package.Type}}Option{ {{toLower $option.Package.Type}}: {{toLower $option.Package.Type}} } }
{{- end}}
{{- end}}
{{end -}}
`)
}
