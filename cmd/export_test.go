package cmd

import (
	"testing"
)

func TestExport(t *testing.T) {
	ExportTest(t, func() {
		given, when, then := NewExportStage(t)

		given.
			spotify_returns_one_show().and().
			itunes_returns_a_match()

		when.
			the_command_is_run()

		then.
			the_output_opml_file_exists().and().
			the_output_opml_file_contains_the_expected_show()
	})
}
