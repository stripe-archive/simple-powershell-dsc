package memory

import (
	"context"
	"sync"

	"github.com/stripe-archive/simple-powershell-dsc/dsc"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/status"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type NodeStatus struct {
	regs     map[string][]string
	regsLock sync.RWMutex
	config   dsc.ConfigurationRepository
}

func New(config dsc.ConfigurationRepository) *NodeStatus {
	return &NodeStatus{
		regs:   make(map[string][]string),
		config: config,
	}
}

func (s *NodeStatus) RegisterDscAgent(
	ctx context.Context,
	req types.RegisterDscAgentRequest,
) (*types.RegisterDscAgentResponse, error) {
	// Save the list of configuration names for this node.
	s.regsLock.Lock()
	defer s.regsLock.Unlock()
	s.regs[req.AgentID] = req.Body.ConfigurationNames

	// No response needed
	return nil, nil
}

func (s *NodeStatus) GetDscAction(
	ctx context.Context,
	req types.GetDscActionRequest,
) (*types.GetDscActionResponse, error) {
	s.regsLock.RLock()
	regs, ok := s.regs[req.AgentID]
	s.regsLock.RUnlock()

	// Validate we've seen a registration
	if !ok {
		return nil, types.AgentNotRegisteredError{req.AgentID}
	}

	return status.ReconcileDscStatus(ctx, s.config, regs, &req)
}
