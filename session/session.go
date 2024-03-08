package session

import (
	"sync"

	"github.com/andaru/gnotch/device"
	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	DefaultManagerSessionCacheSize = 100
)

// Session manages an individual device session
type Session struct {
	device.Device
}

// Manager is a manager of Sessions
type Manager struct {
	mu   sync.Mutex
	size int

	sessions *lru.Cache[string, *Session]
}

// ManagerOption is an option function modifying Manager
type ManagerOption func(*Manager)

// WithSessions is an option setting the session cache size of a Manager.
// size must be a positive integer, else NewManager will panic.
func WithSessions(size int) ManagerOption { return func(m *Manager) { m.size = size } }

// NewManager returns a new session configured with options
func NewManager(options ...ManagerOption) *Manager {
	m := &Manager{}
	for _, option := range options {
		option(m)
	}
	if m.size < 1 {
		m.size = DefaultManagerSessionCacheSize
	}
	m.sessions, _ = lru.NewWithEvict[string, *Session](m.size, m.onEvicted)
	return m
}

func (m *Manager) Session(name string) (s *Session, err error) {
	var ok bool

	defer m.mu.Unlock()
	m.mu.Lock()

	s, ok = m.sessions.Get(name)
	if !ok {
		s, err = m.newSession(name)
		_ = m.sessions.Add(name, s)
		return s, err
	}
	return s, nil
}

func (m *Manager) onEvicted(_ string, v *Session) { v.Device.Close() }

func (m *Manager) newSession(name string) (*Session, error) {
	return nil, nil
}
