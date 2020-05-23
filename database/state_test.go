package database

import "testing"

func TestState(t *testing.T) {
	t.Run("NewStateFromDisk should return a new State struct", func(t *testing.T) {
		if 2 == 2 {
			t.Error("Red")
		}
	})
}
