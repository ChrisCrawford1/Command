package responses

import "testing"

func TestMessageForTag(t *testing.T) {
	t.Run("Will return required tag message", func(t *testing.T) {
		tagValue := MessageForTag("required")

		if tagValue != "Field is required" {
			t.Errorf("Incorrect message for required tag")
		}
	})

	t.Run("Will return default tag message if not specified", func(t *testing.T) {
		tagValue := MessageForTag("")

		if tagValue != "" {
			t.Errorf("Incorrect message for default")
		}
	})
}
