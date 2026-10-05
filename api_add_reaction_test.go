package zca

import (
	"context"
	"testing"
)

func TestAddReactionInvalidThreadType(t *testing.T) {
	_, err := (&API{}).AddReaction(context.Background(), ReactionsLIKE, AddReactionDestination{Type: ThreadType(5)})
	if err == nil || err.Error() != "Thread type is invalid" {
		t.Fatalf("got %v", err)
	}
}
