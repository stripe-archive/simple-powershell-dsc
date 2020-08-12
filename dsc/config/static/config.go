package static

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type ModuleSpec struct {
	Name    string
	Version string
}

type module struct {
	content []byte
	hash    string
}

type ConfigurationRepository struct {
	config     []byte
	configHash string

	modules map[ModuleSpec]module
}

func New(config []byte, modules map[ModuleSpec][]byte) *ConfigurationRepository {
	mods := make(map[ModuleSpec]module)
	if modules != nil {
		// From section 3.7.5:
		//     The server MUST use case-insensitive ordinal
		//     comparison to match ModuleName and ModuleVersion.
		//
		// Thus, we lower-case the module name and version when storing.
		for k, v := range modules {
			spec := ModuleSpec{
				Name:    strings.ToLower(k.Name),
				Version: strings.ToLower(k.Version),
			}
			h := sha256.Sum256(v)
			mods[spec] = module{
				content: v,
				hash:    strings.ToUpper(hex.EncodeToString(h[:])),
			}
		}
	}

	configHash := sha256.Sum256(config)
	return &ConfigurationRepository{
		config:     config,
		configHash: strings.ToUpper(hex.EncodeToString(configHash[:])),
		modules:    mods,
	}
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
	ret := &types.GetConfigurationResponse{
		//ConfigurationName: req.ConfigurationName,
		Content:           bytes.NewReader(c.config),
		Checksum:          c.configHash,
		ChecksumAlgorithm: "SHA-256",
	}
	return ret, nil
}

func (c *ConfigurationRepository) GetConfigurationHash(
	ctx context.Context,
	req types.GetConfigurationRequest,
) (hash, algo string, err error) {
	return c.configHash, "SHA-256", nil
}

func (c *ConfigurationRepository) GetModule(
	ctx context.Context,
	req types.GetModuleRequest,
) (*types.GetModuleResponse, error) {
	// Try to find this module. As above, we must use case-insensitive
	// comparison, so we lowercase before searching.
	spec := ModuleSpec{
		Name:    strings.ToLower(req.Name),
		Version: strings.ToLower(req.Version),
	}
	mod, ok := c.modules[spec]
	if !ok {
		return nil, types.ModuleNotFoundError{
			AgentID: req.AgentID,
			Name:    req.Name,
			Version: req.Version,
		}
	}

	// All set!
	ret := &types.GetModuleResponse{
		Content:           bytes.NewReader(mod.content),
		Checksum:          mod.hash,
		ChecksumAlgorithm: "SHA-256",
	}
	return ret, nil
}
