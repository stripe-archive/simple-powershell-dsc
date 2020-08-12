package local

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/atomic"

	"github.com/stripe-archive/simple-powershell-dsc/dsc"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/status"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type NodeStatus struct {
	path   string
	config dsc.ConfigurationRepository
}

func New(config dsc.ConfigurationRepository, path string) (*NodeStatus, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !st.IsDir() {
		return nil, fmt.Errorf("dsc/status/local: path is not a directory")
	}

	ret := &NodeStatus{
		path:   path,
		config: config,
	}
	return ret, nil
}

func (s *NodeStatus) RegisterDscAgent(
	ctx context.Context,
	req types.RegisterDscAgentRequest,
) (*types.RegisterDscAgentResponse, error) {
	// Atomically write registrations to file.
	body, err := json.Marshal(&req.Body.ConfigurationNames)
	if err != nil {
		return nil, err
	}

	err = atomic.WriteFile(filepath.Join(s.path, req.AgentID+".json"), bytes.NewReader(body))
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
	// Read the file from on disk
	f, err := os.Open(filepath.Join(s.path, req.AgentID+".json"))
	if err != nil {
		// If the file doesn't exist, the agent isn't registered
		if os.IsNotExist(err) {
			return nil, types.AgentNotRegisteredError{req.AgentID}
		}

		return nil, err
	}
	defer f.Close()

	var regs []string
	if err = json.NewDecoder(f).Decode(&regs); err != nil {
		return nil, err
	}

	return status.ReconcileDscStatus(ctx, s.config, regs, &req)
}
