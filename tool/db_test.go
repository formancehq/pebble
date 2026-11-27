// Copyright 2019 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package tool

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cockroachdb/pebble/v2/vfs"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	runTests(t, "testdata/db_*")
}

func TestDBCommandErrors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{name: "check missing database", args: []string{"db", "check", "missing"}, want: "missing"},
		{name: "get missing database", args: []string{"db", "get", "missing", "key"}, want: "missing"},
		{name: "io bench invalid sizes", args: []string{"db", "io-bench", "missing", "--io-sizes-kb=invalid"}, want: "error parsing io-sizes"},
		{name: "analyze data missing directory", args: []string{"db", "analyze-data", "missing", "--output=results.csv"}, want: "loading file list"},
		{name: "wal dump missing file", args: []string{"wal", "dump", "missing.log"}, want: "missing.log"},
		{name: "wal merged invalid filename", args: []string{"wal", "dump-merged", "not-a-log"}, want: "does not parse as a log file"},
		{name: "remote catalog missing file", args: []string{"remotecat", "dump", "missing"}, want: "missing"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := &cobra.Command{SilenceUsage: true, SilenceErrors: true}
			root.AddCommand(New(FS(vfs.NewMem())).Commands...)
			root.SetArgs(tc.args)
			var output bytes.Buffer
			root.SetOut(&output)
			root.SetErr(&output)
			err := root.Execute()
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), tc.want), "error: %v", err)
		})
	}
}
