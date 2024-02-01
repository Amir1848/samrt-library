package library

import (
	"context"
	"testing"
)

func TestLib(t *testing.T) {
	err := SetLibItemStatus(context.Background(), 1, "A1", 1, "fum")
	if err != nil {
		t.Fatal(err.Error())
	}
}
