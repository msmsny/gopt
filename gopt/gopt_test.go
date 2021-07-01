package gopt

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	testCases := map[string]struct {
		name          string
		options       string
		packageName   string
		evaluate      bool
		header        bool
		fileName      string
		formatImports bool
	}{
		"basic": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool",
			packageName: "",
			evaluate:    true,
			fileName:    "testdata/basic.go",
		},
		"withStringSlice": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:stringSlice",
			packageName: "testdata",
			evaluate:    true,
			fileName:    "testdata/withstringslice.go",
		},
		"withPackage": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool",
			packageName: "testdata",
			evaluate:    true,
			fileName:    "testdata/withpackage.go",
		},
		"withoutEvaluate": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool",
			packageName: "testdata",
			evaluate:    false,
			fileName:    "testdata/withoutevaluate.go",
		},
		"format": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool,qux:int64,quux:string",
			packageName: "testdata",
			evaluate:    true,
			fileName:    "testdata/format.go",
		},
		"withDuration": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:duration",
			packageName: "testdata",
			evaluate:    true,
			fileName:    "testdata/withduration.go",
		},
		"withPackageType": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:*go.uber.org/zap.Logger",
			packageName: "testdata",
			evaluate:    true,
			fileName:    "testdata/withpackagetype.go",
		},
		"formatImports": {
			name:          "sample",
			options:       "foo:string,bar:duration,baz:*go.uber.org/zap.Logger",
			packageName:   "testdata",
			evaluate:      true,
			fileName:      "testdata/formatimports.go",
			formatImports: true,
		},
		"withLocalPackageType": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:*LocalPackageType",
			packageName: "testdata",
			evaluate:    true,
			fileName:    "testdata/withlocalpackagetype.go",
		},
		"withHeader": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool",
			packageName: "testdata",
			evaluate:    true,
			header:      true,
			fileName:    "testdata/withheader.go",
		},
	}

	for testCase, tt := range testCases {
		t.Run(testCase, func(t *testing.T) {
			cmd := NewGoptCommand()
			destination := filepath.FromSlash(t.TempDir() + "/" + tt.fileName)
			require.NoError(t, cmd.Flags().Set("name", tt.name))
			require.NoError(t, cmd.Flags().Set("options", tt.options))
			require.NoError(t, cmd.Flags().Set("package", tt.packageName))
			require.NoError(t, cmd.Flags().Set("output", destination))
			if !tt.evaluate {
				require.NoError(t, cmd.Flags().Set("evaluate", "false"))
			}
			if !tt.header {
				require.NoError(t, cmd.Flags().Set("header", "false"))
			}
			if tt.formatImports {
				require.NoError(t, cmd.Flags().Set("format-imports", "true"))
			}
			require.NoError(t, cmd.Execute())

			want, err := ioutil.ReadFile(tt.fileName)
			require.NoError(t, err)
			got, err := ioutil.ReadFile(destination)
			require.NoError(t, err)
			assert.Equal(t, string(want), string(got))
		})
	}
}
