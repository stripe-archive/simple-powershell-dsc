package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type ReportServer struct {
	bucket *string
	s3     s3iface.S3API
}

func New(bucket string, s3 s3iface.S3API) *ReportServer {
	return &ReportServer{&bucket, s3}
}

func (c *ReportServer) RegisterDscAgent(
	ctx context.Context,
	req types.RegisterDscAgentRequest,
) (*types.RegisterDscAgentResponse, error) {
	// No registration required
	return nil, nil
}

func (c *ReportServer) SendReport(
	ctx context.Context,
	req types.SendReportRequest,
) (*types.SendReportResponse, error) {
	key := fmt.Sprintf("reports/%s/%s.json",
		strings.ToLower(req.AgentID),
		strings.ToLower(req.Body.JobID),
	)

	body, err := json.Marshal(&req.Body)
	if err != nil {
		return nil, err
	}

	_, err = c.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      c.bucket,
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ACL:         aws.String("private"),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}

	return &types.SendReportResponse{}, nil
}

func (c *ReportServer) GetReports(
	ctx context.Context,
	req types.GetReportsRequest,
) (*types.GetReportsResponse, error) {
	key := fmt.Sprintf("reports/%s/%s.json",
		strings.ToLower(req.AgentID),
		strings.ToLower(req.JobID),
	)

	result, err := c.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: c.bucket,
		Key:    &key,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, types.ReportNotFoundError{
					AgentID: req.AgentID,
					JobID:   req.JobID,
				}
			}
		}

		// Unknown error; return as-is
		return nil, err
	}
	defer result.Body.Close()

	// Read report body into memory
	report, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return &types.GetReportsResponse{Response: report}, nil
}
