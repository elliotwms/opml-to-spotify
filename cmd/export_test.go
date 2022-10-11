package cmd

import (
	"testing"
)

func TestExport(t *testing.T) {
	ExportTest(t, func() {
		given, when, then := NewExportStage(t)

		given.
			spotify_will_return_one_show().and().
			itunes_will_return_a_match()

		when.
			the_command_is_run()

		then.
			the_output_opml_file_exists().and().
			the_output_opml_file_contains_the_expected_show().and().
			no_errors_are_output()
	})
}
