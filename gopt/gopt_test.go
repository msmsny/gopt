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
		name        string
		options     string
		packageName string
		evaluate    string
		fileName    string
	}{
		"basic": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool",
			packageName: "",
			evaluate:    "true",
			fileName:    "testdata/basic.go",
		},
		"withPackage": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool",
			packageName: "testdata",
			evaluate:    "true",
			fileName:    "testdata/withpackage.go",
		},
		"withoutEvaluate": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool",
			packageName: "testdata",
			evaluate:    "false",
			fileName:    "testdata/withoutevaluate.go",
		},
		"format": {
			name:        "sample",
			options:     "foo:string,bar:int,baz:bool,qux:int64,quux:string",
			packageName: "testdata",
			evaluate:    "true",
			fileName:    "testdata/format.go",
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
			require.NoError(t, cmd.Flags().Set("evaluate", tt.evaluate))
			require.NoError(t, cmd.Execute())

			want, err := ioutil.ReadFile(tt.fileName)
			require.NoError(t, err)
			got, err := ioutil.ReadFile(destination)
			require.NoError(t, err)
			assert.Equal(t, string(want), string(got))
		})
	}
}
