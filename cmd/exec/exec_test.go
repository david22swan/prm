package exec_test

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/puppetlabs/prm/cmd/exec"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func nullFunction(cmd *cobra.Command, args []string) error {
	return nil
}

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		returnCode int
		out        string
		wantCmd    *cobra.Command
		wantErr    bool
		f          func(cmd *cobra.Command, args []string) error
	}{
		{
			name:    "executes without error",
			f:       nullFunction,
			out:     "",
			wantErr: false,
		},
		{
			name:    "executes without error for valid flag",
			args:    []string{"author/templateId"},
			f:       nullFunction,
			out:     "",
			wantErr: false,
		},
		{
			name:    "executes with error when tool provided in the wrong format",
			args:    []string{"foo-bar"},
			f:       nullFunction,
			out:     "Selected tool must be in AUTHOR/ID format",
			wantErr: true,
		},
		{
			name:    "executes with error for invalid flag",
			args:    []string{"--foo"},
			f:       nullFunction,
			out:     "unknown flag: --foo",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			prmObj := &prm.Prm{
				AFS:  &afero.Afero{Fs: fs},
				IOFS: &afero.IOFS{Fs: fs},
			}
			cmd := exec.CreateCommand(prmObj)
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)
			cmd.RunE = tt.f

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("executeTestUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			out, err := ioutil.ReadAll(b)
			if err != nil {
				t.Errorf("Failed to read stdout: %v", err)
				return
			}

			output := string(out)
			r := regexp.MustCompile(tt.out)
			if !r.MatchString(output) {
				t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
				return
			}
		})
	}
}
