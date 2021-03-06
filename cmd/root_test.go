package cmd

import (
	"github.com/freshautomations/sconfig/defaults"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckArgs(t *testing.T) {
	cmd := &cobra.Command{
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Version: defaults.Version,
	}

	assert.NotNil(t, CheckArgs(cmd, []string{"../test.toml"}), "enough parameters")
	assert.NotNil(t, CheckArgs(cmd, []string{"notexist.toml", "myfakeint"}), "file found")
	assert.Nil(t, CheckArgs(cmd, []string{"../test.toml", "myfakeint=3"}), "parameter check")
}

func TestRunRoot(t *testing.T) {
/* Todo: Fix test cases
	cmd := &cobra.Command{
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Version: defaults.Version,
	}

	var result interface{}
	var err error

	// Root section
	result, err = RunRoot(cmd, []string{"../test.toml", "myfakeint=4"})
	assert.Equal(t, "4", result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// Section
	result, err = RunRoot(cmd, []string{"../test.ini", "district9.name"})
	assert.Equal(t, "Wikus", result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// Case-insensitive
	result, err = RunRoot(cmd, []string{"../test.ini", "district9.eta"})
	assert.Equal(t, "3", result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// No key
	result, err = RunRoot(cmd, []string{"../test.ini", "blur2.song3"})
	assert.Empty(t, result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// No section
	result, err = RunRoot(cmd, []string{"../test.ini", "blur3.song3"})
	assert.Empty(t, result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// invalid input
	result, err = RunRoot(cmd, []string{"../test.ini", "#invalid"})
	assert.Empty(t, result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// slices with numbers
	result, err = RunRoot(cmd, []string{"../test.ini", "slices.numbers"})
	assert.Equal(t, "1 1 2 3 5", result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// slices with numbers and strings
	result, err = RunRoot(cmd, []string{"../test.ini", "slices.strings"})
	assert.Equal(t, "The answer is 42", result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// reading keys of a section
	result, err = RunRoot(cmd, []string{"../test.ini", "slices"})
	assert.Equal(t, "numbers strings", result, "unexpected result")
	assert.Nil(t, err, "unexpected error")

	// reading section names (keys of the root section)
	result, err = RunRoot(cmd, []string{"../test.ini", "."})
	assert.Equal(t, "blur district9 master_of_the_universe slices", result, "unexpected result")
	assert.Nil(t, err, "unexpected error")
*/
}
