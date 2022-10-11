package cmd

import "testing"

func TestImport(t *testing.T) {
	ImportTest(t, func() {
		given, when, then := NewImportStage(t)

		given.
			an_opml_file().and().
			spotify_will_return_search_results().and().
			spotify_will_all_the_user_to_save_the_shows()

		when.
			the_command_is_run()

		then.
			the_user_is_subscribed_to_the_show()
	})
}
