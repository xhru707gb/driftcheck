package tfstate

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteWatchlistReport_Nil(t *testing.T) {
	var buf bytes.Buffer
	WriteWatchlistReport(&buf, nil)
	if !strings.Contains(buf.String(), "nil") {
		t.Errorf("expected nil in output, got: %s", buf.String())
	}
}

func TestWriteWatchlistReport_MultipleEntries(t *testing.T) {
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "web"})
	wl.Add(WatchEntry{ResourceType: "aws_s3_bucket", ResourceName: "data",
		Attributes: []string{"bucket"}})

	var buf bytes.Buffer
	WriteWatchlistReport(&buf, wl)
	out := buf.String()

	if !strings.Contains(out, "2 resource(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "aws_s3_bucket.data") {
		t.Errorf("expected s3 bucket in output, got: %s", out)
	}
	if !strings.Contains(out, "<all attributes>") {
		t.Errorf("expected all-attributes marker for aws_instance.web, got: %s", out)
	}
}

func TestWatchlist_String_Empty(t *testing.T) {
	wl := NewWatchlist()
	s := wl.String()
	if !strings.Contains(s, "all resources") {
		t.Errorf("unexpected string for empty watchlist: %s", s)
	}
}

func TestWatchlist_String_WithEntries(t *testing.T) {
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "app",
		Attributes: []string{"ami"}})
	s := wl.String()
	if !strings.Contains(s, "aws_instance.app") {
		t.Errorf("expected resource in string, got: %s", s)
	}
	if !strings.Contains(s, "ami") {
		t.Errorf("expected attribute in string, got: %s", s)
	}
}
