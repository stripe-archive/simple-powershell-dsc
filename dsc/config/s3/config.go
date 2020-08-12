package s3

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type ConfigurationRepository struct {
	bucket *string
	s3     s3iface.S3API
}

func New(bucket string, s3 s3iface.S3API) *ConfigurationRepository {
	return &ConfigurationRepository{&bucket, s3}
}

func (c *ConfigurationRepository) RegisterDscAgent(
	ctx context.Context,
	req types.RegisterDscAgentRequest,
) (*types.RegisterDscAgentResponse, error) {
	// No registration required
	return nil, nil
}

func (c *ConfigurationRepository) GetConfiguration(
	ctx context.Context,
	req types.GetConfigurationRequest,
) (*types.GetConfigurationResponse, error) {
	// From section 3.6.5:
	//     The server MUST use case-insensitive ordinal comparison to match
	//     the AgentId and ConfigurationName.
	//
	// Thus, we lower-case the module name and version, since S3 doesn't
	// support case-insensitive comparison.
	key := fmt.Sprintf("config/%s",
		strings.ToLower(req.ConfigurationName),
	)

	// Attempt to fetch the configuration from S3
	input := &s3.GetObjectInput{
		Bucket: c.bucket,
		Key:    &key,
	}

	result, err := c.s3.GetObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, types.ConfigurationNotFoundError{
					AgentID: req.AgentID,
					Name:    req.ConfigurationName,
				}
			}
		}

		// Unknown error; return as-is
		return nil, err
	}
	defer result.Body.Close()

	// Read config body into memory
	config, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	// Hash body
	h := sha256.Sum256(config)

	// All set!
	ret := &types.GetConfigurationResponse{
		//ConfigurationName: req.ConfigurationName,
		Content:           bytes.NewReader(config),
		Checksum:          strings.ToUpper(hex.EncodeToString(h[:])),
		ChecksumAlgorithm: "SHA-256",
	}
	return ret, nil
}

func (c *ConfigurationRepository) GetModule(
	ctx context.Context,
	req types.GetModuleRequest,
) (*types.GetModuleResponse, error) {
	// From section 3.7.5:
	//     The server MUST use case-insensitive ordinal
	//     comparison to match ModuleName and ModuleVersion.
	//
	// Thus, we lower-case the module name and version, since S3 doesn't
	// support case-insensitive comparison.
	key := fmt.Sprintf("modules/%s/%s.zip",
		strings.ToLower(req.Name),
		strings.ToLower(req.Version),
	)

	// Attempt to fetch the configuration from S3
	input := &s3.GetObjectInput{
		Bucket: c.bucket,
		Key:    &key,
	}

	result, err := c.s3.GetObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, types.ModuleNotFoundError{
					AgentID: req.AgentID,
					Name:    req.Name,
					Version: req.Version,
				}
			}
		}

		// Unknown error; return as-is
		return nil, err
	}
	defer result.Body.Close()

	// Read config body into memory
	config, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	// Hash body
	h := sha256.Sum256(config)

	// All set!
	ret := &types.GetModuleResponse{
		Content:           bytes.NewReader(config),
		Checksum:          strings.ToUpper(hex.EncodeToString(h[:])),
		ChecksumAlgorithm: "SHA-256",
	}
	return ret, nil
}
