package cmd

import "testing"

func TestImport(t *testing.T) {
	ImportTest(t, func() {
		given, when, then := NewImportStage(t)

		given.
			an_opml_file().and().
			spotify_will_return_search_results().and().
			spotify_will_allow_the_user_to_save_the_shows()

		when.
			the_command_is_run()

		then.
			the_user_is_subscribed_to_the_show().and().
			no_errors_are_output()
	})
}

func TestImport_NotFound_NoResults(t *testing.T) {
	ImportTest(t, func() {
		given, when, then := NewImportStage(t)

		given.
			an_opml_file().and().
			spotify_will_return_no_search_results()

		when.
			the_command_is_run()

		then.
			the_user_is_not_subscribed_to_any_shows().and().
			the_error_is_output("Could not find show: Hello, World!")
	})
}

func TestImport_NotFound_NoExactMatch(t *testing.T) {
	ImportTest(t, func() {
		given, when, then := NewImportStage(t)

		given.
			an_opml_file().and().
			spotify_will_return_no_exact_matching_results()

		when.
			the_command_is_run()

		then.
			the_user_is_not_subscribed_to_any_shows().and().
			the_error_is_output("Could not find show: Hello, World!")
	})
}

func TestImport_DryRun(t *testing.T) {
	ImportTest(t, func() {
		given, when, then := NewImportStage(t)

		given.
			an_opml_file().and().
			spotify_will_return_search_results().and().
			the_dry_run_flag_is_set()

		when.
			the_command_is_run()

		then.
			the_user_is_not_subscribed_to_any_shows().and().
			no_errors_are_output().and().
			the_message_is_output("Dry-run. Exiting...")
	})
}

func TestImport_WithMissingFlag(t *testing.T) {
	ImportTest(t, func() {
		given, when, then := NewImportStage(t)

		given.
			an_opml_file().and().
			spotify_will_return_no_search_results().and().
			the_missing_flag_is_set()

		when.
			the_command_is_run()

		then.
			the_user_is_not_subscribed_to_any_shows().and().
			the_message_is_output("Writing 1 missing show titles to missing.txt").and().
			the_missing_file_exists()
	})
}
