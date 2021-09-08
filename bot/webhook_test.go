package bot

import (
	"bufio"
	"strings"
	"testing"
)

func TestDecodeWebhook(t *testing.T) {
	body := strings.NewReader(`{"object_kind": "merge_request"}`)
	got, _ := decodeWebhook(bufio.NewReader(body))
	if got.ObjectKind != MergeRequest {
		t.Errorf("expected %s, but got: %v", MergeRequest, got.ObjectKind)
	}
}
