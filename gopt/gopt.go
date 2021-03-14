package gopt

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func NewGoptCommand() *cobra.Command {
	var (
		name        *string
		options     *[]string
		packageName *string
		destination *string
		evaluate    *bool
		flagErrors  []error
	)
	cmds := &cobra.Command{
		Use:           "gopt",
		Short:         "gopt generates functional options pattern code",
		Long:          "gopt generates functional options pattern code",
		SilenceErrors: true,
		SilenceUsage:  false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			for _, err := range flagErrors {
				if err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			gopt := &gopt{
				tpl: template.New("gopt").Funcs(map[string]interface{}{
					"toLower": func(name string) string {
						return strings.ToLower(name)
					},
					"title": func(name string) string {
						return strings.Title(name)
					},
					"hasImport": func(opts []*templateOption) bool {
						for _, opt := range opts {
							if opt.IsDuration() || opt.IsPackage() {
								return true
							}
						}

						return false
					},
				}),
				params: &templateParams{
					Name:        *name,
					PackageName: *packageName,
					Evaluate:    *evaluate,
				},
				writer: os.Stdout,
				dest:   *destination,
			}
			tplOpts, err := parseOptions(*options)
			if err != nil {
				return fmt.Errorf("parseOptions: %s", err)
			}
			gopt.params.Options = tplOpts

			return gopt.run()
		},
	}

	flags := cmds.Flags()
	flags.SortFlags = false
	name = flags.String("name", "", "functional options name to specify variadic functions arguments (required)")
	flagErrors = append(flagErrors, cobra.MarkFlagRequired(flags, "name"))
	options = flags.StringSlice("options", []string{}, "option names and values, e.g.: foo:string,bar:int,baz:bool")
	flagErrors = append(flagErrors, cobra.MarkFlagRequired(flags, "options"))
	packageName = flags.String("package", "", "output package name")
	destination = flags.StringP("output", "o", "", "output file name")
	evaluate = flags.Bool("evaluate", true, "output evaluateOptions")

	return cmds
}

type templateParams struct {
	Name        string
	Options     []*templateOption
	PackageName string
	Evaluate    bool
}

type templateOption struct {
	Name    string
	Type    string
	Package *templatePackage
}

func (o *templateOption) IsDuration() bool {
	return o.Type == "time.Duration"
}

func (o *templateOption) IsPackage() bool {
	return o.Package != nil
}

type templatePackage struct {
	Prefix string // []*
	Base   string // path/to
	Name   string // repo
	Type   string // MyType
}

func (i *templatePackage) IsFullPath() bool {
	return i.Base != "" && i.Name != ""
}

func (i *templatePackage) HasName() bool {
	return i.Base == "" && i.Name != ""
}

type gopt struct {
	params *templateParams
	tpl    *template.Template
	writer io.Writer
	dest   string
}

func (g *gopt) run() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %s", err)
	}

	rawTpl := uberGoStyleOptionTemplate()
	t, err := g.tpl.Parse(string(rawTpl))
	if err != nil {
		return fmt.Errorf("tpl.Parse: %s", err)
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, g.params); err != nil {
		return fmt.Errorf("t.Execute: %s", err)
	}
	contents, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format.Source: %s", err)
	}

	if g.dest != "" {
		path := filepath.FromSlash(g.dest)
		if !filepath.IsAbs(path) {
			path = filepath.FromSlash(wd + "/" + g.dest)
		}
		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("os.MkdirAll: %s", err)
				}
			} else {
				return fmt.Errorf("os.Stat: %s", err)
			}
		}
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("os.Create: %s", err)
		}
		defer file.Close()
		g.writer = file
	}

	if _, err := g.writer.Write(contents); err != nil {
		return fmt.Errorf("g.writer.Write: %s", err)
	}

	return nil
}

// prefix, path(domain), name, type
var packageTypeRegexp = regexp.MustCompile(`^(?:((?:\[\])?\*?)(?:(.+)\/)?([a-zA-Z0-9]+)(?:\.))?([a-zA-Z0-9]+)$`)

// parseOptions parse names and types like "foo:string,bar:int,baz:bool"
func parseOptions(opts []string) ([]*templateOption, error) {
	var tplOpts []*templateOption
	for _, opt := range opts {
		nameType := strings.Split(opt, ":")
		if len(nameType) != 2 {
			return nil, fmt.Errorf("invalid options format: %s", opt)
		}
		switch nameType[1] {
		// TODO cover
		//  int  int8  int16  int32  int64
		//  uint uint8 uint16 uint32 uint64 uintptr
		//  byte // alias for uint8
		//  rune // alias for int32
		//  float32 float64
		//  complex64 complex128
		case "string", "int", "int64", "bool", "duration", "stringSlice":
			tplOpt := &templateOption{
				Name: nameType[0],
				Type: nameType[1],
			}
			if nameType[1] == "duration" {
				tplOpt.Type = "time.Duration"
			}
			tplOpts = append(tplOpts, tplOpt)
		default:
			packageType := packageTypeRegexp.FindStringSubmatch(nameType[1])
			if packageType == nil {
				return nil, fmt.Errorf("option type %s must be string, int, int64, bool, duration, stringSlice or custom type", nameType[1])
			}
			if len(packageType) != 5 {
				return nil, fmt.Errorf("custom type must be 5 columns of original value, prefix, base path, name, type")
			}
			tplOpt := &templateOption{
				Name: nameType[0],
				Package: &templatePackage{
					Prefix: packageType[1],
					Base:   packageType[2],
					Name:   packageType[3],
					Type:   packageType[4],
				},
			}
			tplOpts = append(tplOpts, tplOpt)
		}
	}

	return tplOpts, nil
}
