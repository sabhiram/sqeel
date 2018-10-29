package sqeel

import "testing"

func TestSnakeCase(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected string
	}{
		{"HelloThere", "hello_there"},
		{"ID", "id"},
		{"SweetIDThatIsAwesome", "sweet_id_that_is_awesome"},
	} {
		act := ToSnakeCase(tc.input)
		if act != tc.expected {
			t.Fatalf("Error, expected (%s) got (%s)\n", tc.expected, act)
		}
	}
}
