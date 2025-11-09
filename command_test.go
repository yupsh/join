package command_test

import (
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/join"
)

func TestJoin_Basic(t *testing.T) {
	result := run.Quick(command.Join("testdata/names.txt", "testdata/scores.txt"))
	assertion.NoError(t, result.Err)
	// Should join on first field (IDs 1, 2, 3 are in both files)
	assertion.Count(t, result.Stdout, 3)
}

func TestJoin_Field1(t *testing.T) {
	result := run.Quick(command.Join("testdata/names.txt", "testdata/scores.txt", command.Field1(1)))
	assertion.NoError(t, result.Err)
	// Explicit field 1 (same as default)
}

func TestJoin_Field2(t *testing.T) {
	result := run.Quick(command.Join("testdata/names.txt", "testdata/scores.txt", command.Field2(1)))
	assertion.NoError(t, result.Err)
	// Explicit field 2 (same as default)
}

func TestJoin_Unpaired1(t *testing.T) {
	// Show unpaired lines from file 1 (David with ID 4)
	result := run.Quick(command.Join("testdata/names.txt", "testdata/scores.txt", command.UnpairedFile1))
	assertion.NoError(t, result.Err)
	// Should include unpaired lines from file 1
}

func TestJoin_MissingFile(t *testing.T) {
	result := run.Quick(command.Join("nonexistent.txt", "testdata/scores.txt"))
	assertion.Error(t, result.Err)
}

