package sendshit

import (
	"testing"
)

func TestRandomness(t *testing.T) {
	key, err := generateRandomString(24)

	if err != nil {
		t.Fatal(err)
	}

	if len(key) != 48 {
		t.Error("string not generated")
	}
}
