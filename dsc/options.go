package dsc

import (
	"github.com/sirupsen/logrus"
)

// Option is the type of functional options that can be passed to the
// NewManager function.
type Option func(*Manager)

// WithLogger sets the Logrus logger to use for messages.
func WithLogger(l logrus.FieldLogger) Option {
	return func(m *Manager) {
		m.log = l
	}
}

// WithKeys sets the authentication keys to use to validate incoming requests;
// if none are provided, then no authentication is performed.
//
// TODO(andrew): switch default to "fail auth" if none provided/not insecure
func WithKeys(keys []string) Option {
	return func(m *Manager) {
		m.keys = keys
	}
}
