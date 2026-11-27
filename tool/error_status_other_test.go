package tool

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
)

func TestIntrospectionCommandsReturnMissingFileErrors(t *testing.T) {
	t.Parallel()
	tool := New()
	missing := t.TempDir() + "/missing"
	tests := []struct {
		name    string
		command *cobra.Command
		args    []string
	}{
		{"blob layout", tool.blob.Layout, []string{missing + ".blob"}},
		{"sstable check", tool.sstable.Check, []string{missing + ".sst"}},
		{"sstable layout", tool.sstable.Layout, []string{missing + ".sst"}},
		{"sstable properties", tool.sstable.Properties, []string{missing + ".sst"}},
		{"sstable scan", tool.sstable.Scan, []string{missing + ".sst"}},
		{"sstable space", tool.sstable.Space, []string{missing + ".sst"}},
		{"manifest dump", tool.manifest.Dump, []string{missing}},
		{"manifest summarize", tool.manifest.Summarize, []string{missing}},
		{"manifest check", tool.manifest.Check, []string{missing}},
		{"find", tool.find.Root, []string{missing, "key"}},
		{"lsm", tool.lsm.Root, []string{missing}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.command.SetOut(io.Discard)
			tc.command.SetErr(io.Discard)
			if tc.command.RunE == nil {
				t.Fatal("command has no RunE")
			}
			if err := tc.command.RunE(tc.command, tc.args); err == nil {
				t.Fatal("missing file returned success")
			}
		})
	}
}
