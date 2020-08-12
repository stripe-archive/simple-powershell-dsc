package util

import (
	"context"
	"io"

	"github.com/stripe-archive/simple-powershell-dsc/dsc"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type ConfigRepoHasher interface {
	// GetConfigurationHash returns the hash of a configuration. This
	// method is useful if it's possible to fetch the hash of a
	// configuration without fetching the contents. If this method is not
	// supported, return the empty string and a nil error, and callers
	// should fall back to the `GetConfiguration` method.
	GetConfigurationHash(ctx context.Context, req types.GetConfigurationRequest) (hash, algo string, err error)
}

// GetConfigHash is a helper function that will retrieve the hash of a configuration, attempting to call the `GetConfigurationHash` helper function if it's present.
func GetConfigHash(ctx context.Context, repo dsc.ConfigurationRepository, agentId, configName string) (hash, algo string, err error) {
	req := types.GetConfigurationRequest{
		AgentID:           agentId,
		ConfigurationName: configName,
	}

	if iface, ok := repo.(ConfigRepoHasher); ok {
		// Try to call the helper function
		hash, algo, err := iface.GetConfigurationHash(ctx, req)
		if err == nil && hash != "" {
			return hash, algo, nil
		}

		// Fall through to calling the `GetConfiguration`
	}

	config, err := repo.GetConfiguration(ctx, req)
	if err != nil {
		return "", "", err
	}

	// Ensure we close any response body
	if config.Content != nil {
		if cl, ok := config.Content.(io.Closer); ok {
			cl.Close()
		}
	}

	return config.Checksum, config.ChecksumAlgorithm, nil
}
