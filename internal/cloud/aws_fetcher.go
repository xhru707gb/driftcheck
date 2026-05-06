package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ResourceAttributes holds key-value pairs describing a cloud resource's live state.
type ResourceAttributes map[string]interface{}

// LiveResource represents a single resource fetched from the cloud provider.
type LiveResource struct {
	Type       string
	ID         string
	Attributes ResourceAttributes
}

// Fetcher defines the interface for retrieving live cloud resource state.
type Fetcher interface {
	Fetch(ctx context.Context, resourceType, resourceID string) (*LiveResource, error)
}

// AWSFetcher fetches live resource state from AWS.
type AWSFetcher struct {
	ec2Client *ec2.Client
	s3Client  *s3.Client
}

// NewAWSFetcher creates an AWSFetcher using the default AWS config.
func NewAWSFetcher(ctx context.Context) (*AWSFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}
	return &AWSFetcher{
		ec2Client: ec2.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
	}, nil
}

// Fetch retrieves the live attributes of a resource by type and ID.
func (f *AWSFetcher) Fetch(ctx context.Context, resourceType, resourceID string) (*LiveResource, error) {
	switch resourceType {
	case "aws_instance":
		return f.fetchEC2Instance(ctx, resourceID)
	case "aws_s3_bucket":
		return f.fetchS3Bucket(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (f *AWSFetcher) fetchEC2Instance(ctx context.Context, instanceID string) (*LiveResource, error) {
	out, err := f.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instance %s: %w", instanceID, err)
	}
	if len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}
	inst := out.Reservations[0].Instances[0]
	attrs := ResourceAttributes{
		"instance_type": string(inst.InstanceType),
		"ami":           aws.ToString(inst.ImageId),
		"state":         string(inst.State.Name),
	}
	return &LiveResource{Type: "aws_instance", ID: instanceID, Attributes: attrs}, nil
}

func (f *AWSFetcher) fetchS3Bucket(ctx context.Context, bucketName string) (*LiveResource, error) {
	_, err := f.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("head bucket %s: %w", bucketName, err)
	}
	attrs := ResourceAttributes{
		"bucket": bucketName,
	}
	return &LiveResource{Type: "aws_s3_bucket", ID: bucketName, Attributes: attrs}, nil
}
