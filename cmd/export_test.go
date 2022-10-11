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
			the_output_opml_file_is_created().and().
			the_output_opml_file_contains_the_expected_show().and().
			no_errors_are_output()
	})
}

func TestExport_UserHasNoShows(t *testing.T) {
	given, when, then := NewExportStage(t)

	given.
		spotify_will_return_zero_shows()

	when.
		the_command_is_run()

	then.
		the_output_opml_file_is_not_created().and().
		the_error_is_output("No shows could be found. Exiting")
}

func TestExport_ShowIsNotMatchedIniTunes(t *testing.T) {
	ExportTest(t, func() {
		given, when, then := NewExportStage(t)

		given.
			spotify_will_return_one_show().and().
			itunes_will_not_return_a_match()

		when.
			the_command_is_run()

		then.
			the_output_opml_file_is_not_created().and().
			the_error_is_output("Could not match show: Hello, World!")
	})
}
