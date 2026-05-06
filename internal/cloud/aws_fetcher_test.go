package cloud_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/driftcheck/internal/cloud"
)

// mockFetcher implements Fetcher for testing.
type mockFetcher struct {
	resources map[string]*cloud.LiveResource
	err       error
}

func (m *mockFetcher) Fetch(_ context.Context, resourceType, resourceID string) (*cloud.LiveResource, error) {
	if m.err != nil {
		return nil, m.err
	}
	key := resourceType + "/" + resourceID
	res, ok := m.resources[key]
	if !ok {
		return nil, errors.New("resource not found")
	}
	return res, nil
}

func TestMockFetcher_Found(t *testing.T) {
	f := &mockFetcher{
		resources: map[string]*cloud.LiveResource{
			"aws_instance/i-123": {
				Type: "aws_instance",
				ID:   "i-123",
				Attributes: cloud.ResourceAttributes{
					"instance_type": "t3.micro",
					"ami":           "ami-abc",
				},
			},
		},
	}

	res, err := f.Fetch(context.Background(), "aws_instance", "i-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Attributes["instance_type"] != "t3.micro" {
		t.Errorf("expected t3.micro, got %v", res.Attributes["instance_type"])
	}
}

func TestMockFetcher_NotFound(t *testing.T) {
	f := &mockFetcher{resources: map[string]*cloud.LiveResource{}}
	_, err := f.Fetch(context.Background(), "aws_instance", "i-missing")
	if err == nil {
		t.Fatal("expected error for missing resource")
	}
}

func TestMockFetcher_Error(t *testing.T) {
	f := &mockFetcher{err: errors.New("api error")}
	_, err := f.Fetch(context.Background(), "aws_instance", "i-123")
	if err == nil {
		t.Fatal("expected error from fetcher")
	}
}

func TestLiveResource_Fields(t *testing.T) {
	r := &cloud.LiveResource{
		Type: "aws_s3_bucket",
		ID:   "my-bucket",
		Attributes: cloud.ResourceAttributes{"bucket": "my-bucket"},
	}
	if r.Type != "aws_s3_bucket" {
		t.Errorf("unexpected type: %s", r.Type)
	}
	if r.Attributes["bucket"] != "my-bucket" {
		t.Errorf("unexpected bucket attr: %v", r.Attributes["bucket"])
	}
}
