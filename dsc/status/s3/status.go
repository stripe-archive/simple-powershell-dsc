package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/stripe-archive/simple-powershell-dsc/dsc"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/status"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type NodeStatus struct {
	bucket *string
	s3     s3iface.S3API
	config dsc.ConfigurationRepository
}

func New(config dsc.ConfigurationRepository, bucket string, s3 s3iface.S3API) *NodeStatus {
	return &NodeStatus{
		bucket: &bucket,
		s3:     s3,
		config: config,
	}
}

func (s *NodeStatus) RegisterDscAgent(
	ctx context.Context,
	req types.RegisterDscAgentRequest,
) (*types.RegisterDscAgentResponse, error) {
	key := fmt.Sprintf("registrations/%s.json", strings.ToLower(req.AgentID))

	body, err := json.Marshal(&req.Body)
	if err != nil {
		return nil, err
	}

	_, err = s.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      s.bucket,
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ACL:         aws.String("private"),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}

	// No response needed
	return nil, nil
}

func (s *NodeStatus) GetDscAction(
	ctx context.Context,
	req types.GetDscActionRequest,
) (*types.GetDscActionResponse, error) {
	key := fmt.Sprintf("registrations/%s.json", strings.ToLower(req.AgentID))

	// Get registration from S3
	result, err := s.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    &key,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, types.AgentNotRegisteredError{req.AgentID}
			}
		}

		// Unknown error; return as-is
		return nil, err
	}
	defer result.Body.Close()

	// Deserialize into body
	var registerBody types.RegisterDscAgentRequestBody
	if err := json.NewDecoder(result.Body).Decode(&registerBody); err != nil {
		return nil, err
	}

	return status.ReconcileDscStatus(ctx, s.config, registerBody.ConfigurationNames, &req)
}
