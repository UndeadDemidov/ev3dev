// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ev3dev

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

var mockedExitStatus = 0
var mockedStdout string

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockedExitStatus)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + mockedStdout,
		"EXIT_STATUS=" + es}
	return cmd
}

func TestExecCommandHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	_, _ = fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(i)
}

func TestConsole_readConsoleSize(t *testing.T) {
	tests := []struct {
		name             string
		mockedExitStatus int
		mockedStdout     string
		wantRows         int
		wantCols         int
		wantErr          bool
	}{
		{
			name:             "valid stty size command output",
			mockedExitStatus: 0,
			mockedStdout:     "24 48\n",
			wantRows:         24,
			wantCols:         48,
			wantErr:          false,
		},
	}

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Console{
				columns: 0,
				rows:    0,
			}
			mockedExitStatus = tt.mockedExitStatus
			mockedStdout = tt.mockedStdout

			gotRows, gotCols, err := c.readConsoleSize()
			if (err != nil) != tt.wantErr {
				t.Errorf("readConsoleSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRows != tt.wantRows {
				t.Errorf("readConsoleSize() gotRows = %v, want %v", gotRows, tt.wantRows)
			}
			if gotCols != tt.wantCols {
				t.Errorf("readConsoleSize() gotCols = %v, want %v", gotCols, tt.wantCols)
			}
		})
	}
}
