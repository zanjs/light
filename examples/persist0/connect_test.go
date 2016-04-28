package persist0

import "testing"

func TestConnect(t *testing.T) {
	err := connect()
	if err != nil {
		t.Fatal(err)
	}
}
