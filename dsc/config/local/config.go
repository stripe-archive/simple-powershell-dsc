package local

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/internal/fsmock"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type ConfigurationRepository struct {
	root string
	fs   fsmock.FileSystem
}

func New(root string) *ConfigurationRepository {
	return &ConfigurationRepository{root, fsmock.OSFS{}}
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
	// Thus, we lower-case the module name and version.
	path := filepath.Join(
		c.root,
		"config",
		strings.ToLower(req.ConfigurationName),
	)
	f, err := c.fs.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, types.ConfigurationNotFoundError{
				AgentID: req.AgentID,
				Name:    req.ConfigurationName,
			}
		}

		// Unknown error; return as-is
		return nil, err
	}
	defer f.Close()

	// Read config body into memory
	config, err := ioutil.ReadAll(f)
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
	path := filepath.Join(
		c.root,
		"modules",
		strings.ToLower(req.Name),
		strings.ToLower(req.Version)+".zip",
	)
	f, err := c.fs.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, types.ModuleNotFoundError{
				AgentID: req.AgentID,
				Name:    req.Name,
				Version: req.Version,
			}
		}

		// Unknown error; return as-is
		return nil, err
	}
	defer f.Close()

	config, err := ioutil.ReadAll(f)
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
