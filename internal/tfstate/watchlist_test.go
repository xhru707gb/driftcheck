package tfstate

import (
	"bytes"
	"strings"
	"testing"
)

func buildWatchState() *State {
	s := NewState()
	s.Add(Resource{Type: "aws_instance", Name: "web", ID: "i-001",
		Attributes: map[string]interface{}{"ami": "ami-123", "instance_type": "t2.micro"}})
	s.Add(Resource{Type: "aws_s3_bucket", Name: "data", ID: "bkt-1",
		Attributes: map[string]interface{}{"bucket": "my-bucket", "region": "us-east-1"}})
	s.Add(Resource{Type: "aws_security_group", Name: "sg", ID: "sg-1",
		Attributes: map[string]interface{}{"name": "default"}})
	return s
}

func TestWatchlist_Empty_MatchesAll(t *testing.T) {
	wl := NewWatchlist()
	key := ResourceKey{Type: "aws_instance", Name: "web"}
	if !wl.Matches(key) {
		t.Error("empty watchlist should match all resources")
	}
}

func TestWatchlist_Matches_Present(t *testing.T) {
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "web"})
	if !wl.Matches(ResourceKey{Type: "aws_instance", Name: "web"}) {
		t.Error("expected match")
	}
}

func TestWatchlist_Matches_Absent(t *testing.T) {
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "web"})
	if wl.Matches(ResourceKey{Type: "aws_s3_bucket", Name: "data"}) {
		t.Error("expected no match")
	}
}

func TestWatchlist_WatchedAttributes(t *testing.T) {
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "web",
		Attributes: []string{"ami"}})
	attrs := wl.WatchedAttributes(ResourceKey{Type: "aws_instance", Name: "web"})
	if len(attrs) != 1 || attrs[0] != "ami" {
		t.Errorf("unexpected attrs: %v", attrs)
	}
}

func TestApplyWatchlist_NilState(t *testing.T) {
	wl := NewWatchlist()
	if ApplyWatchlist(nil, wl) != nil {
		t.Error("expected nil for nil state")
	}
}

func TestApplyWatchlist_EmptyWatchlist(t *testing.T) {
	s := buildWatchState()
	wl := NewWatchlist()
	out := ApplyWatchlist(s, wl)
	if out != s {
		t.Error("empty watchlist should return original state")
	}
}

func TestApplyWatchlist_FiltersResources(t *testing.T) {
	s := buildWatchState()
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "web"})
	out := ApplyWatchlist(s, wl)
	keys := out.Keys()
	if len(keys) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(keys))
	}
	if keys[0].Type != "aws_instance" {
		t.Errorf("unexpected resource type: %s", keys[0].Type)
	}
}

func TestApplyWatchlist_FiltersAttributes(t *testing.T) {
	s := buildWatchState()
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "web",
		Attributes: []string{"ami"}})
	out := ApplyWatchlist(s, wl)
	res, _ := out.Get(ResourceKey{Type: "aws_instance", Name: "web"})
	if _, ok := res.Attributes["ami"]; !ok {
		t.Error("expected ami attribute")
	}
	if _, ok := res.Attributes["instance_type"]; ok {
		t.Error("instance_type should have been filtered out")
	}
}

func TestWriteWatchlistReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteWatchlistReport(&buf, NewWatchlist())
	if !strings.Contains(buf.String(), "ALL") {
		t.Errorf("expected ALL in output, got: %s", buf.String())
	}
}

func TestWriteWatchlistReport_WithEntries(t *testing.T) {
	wl := NewWatchlist()
	wl.Add(WatchEntry{ResourceType: "aws_instance", ResourceName: "web",
		Attributes: []string{"ami", "instance_type"}})
	var buf bytes.Buffer
	WriteWatchlistReport(&buf, wl)
	out := buf.String()
	if !strings.Contains(out, "aws_instance.web") {
		t.Errorf("expected resource in report, got: %s", out)
	}
	if !strings.Contains(out, "ami") {
		t.Errorf("expected ami in report, got: %s", out)
	}
}
