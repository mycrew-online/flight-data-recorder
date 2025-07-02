package simconnectmanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	logz "github.com/mrlm-net/go-logz/pkg/logger"
	"github.com/mrlm-net/simconnect/pkg/client"
	"github.com/mrlm-net/simconnect/pkg/types"
	"github.com/mycrew-online/flight-data-recorder/internal/logadapter"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const simStateRequestID uint32 = 1001

// --- SimulatorState for system state monitoring ---
type SimulatorState struct {
	mu      sync.RWMutex
	Sim     int
	Pause   int
	Crashed int
	View    int
}

func (s *SimulatorState) SetSim(val int) {
	s.mu.Lock()
	s.Sim = val
	s.mu.Unlock()
}
func (s *SimulatorState) GetSim() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Sim
}

// ...implement Pause, Crashed, View similarly if needed...

type SimConnectManager struct {
	client   *client.Engine
	state    int
	stateMu  sync.Mutex
	stopCh   chan struct{}
	stopped  sync.WaitGroup
	statusCh chan bool // true=connected, false=disconnected
	logger   *logadapter.LogzWailsAdapter
	simState SimulatorState
	wailsCtx context.Context // Wails context for event emission
}

// SetLogger allows injection of a custom logger (Wails/go-logz adapter)
func (m *SimConnectManager) SetLogger(logger *logadapter.LogzWailsAdapter) {
	m.logger = logger
}

// SetWailsContext sets the Wails context for event emission
func (m *SimConnectManager) SetWailsContext(ctx context.Context) {
	m.wailsCtx = ctx
}

const (
	Offline = iota
	Connecting
	Online
)

// type Logger interface {
//    Info(args ...interface{})
//    Debug(args ...interface{})
//    Warning(args ...interface{})
//    Error(args ...interface{})
//    Fatal(args ...interface{})

func NewSimConnectManager() *SimConnectManager {
	// Create a go-logz logger instance
	lz := logz.NewLogger(logz.LogOptions{
		Level:   logz.Debug,
		Format:  logz.StringOutput,
		Prefix:  "SimConnectManager",
		Outputs: []logz.OutputFunc{logz.ConsoleOutput()},
	})
	// Wrap it with the Wails-compatible adapter
	adapter := logadapter.New(lz)
	return &SimConnectManager{
		stopCh:   make(chan struct{}),
		statusCh: make(chan bool, 1),
		logger:   adapter,
	}
}

// StartConnection starts the connection monitoring goroutine
func (m *SimConnectManager) StartConnection() {
	m.stopCh = make(chan struct{})
	m.stopped.Add(1)
	go func() {
		defer m.stopped.Done()
		retryInterval := 5 * time.Second
		for {
			select {
			case <-m.stopCh:
				m.logDebug("[SimConnectManager] Connection loop stopped.")
				return
			default:
				m.stateMu.Lock()
				if m.state == Offline {
					m.logDebug("[SimConnectManager] State is Offline, will try to connect.")
					m.state = Connecting
					m.stateMu.Unlock()
					m.connect()
				} else {
					m.stateMu.Unlock()
				}
				time.Sleep(retryInterval)
			}
		}
	}()
}

// StopConnection signals the monitoring goroutine to stop and waits for it to finish.
func (m *SimConnectManager) StopConnection() {
	if m.stopCh != nil {
		close(m.stopCh)
		m.stopped.Wait()
		m.disconnect()
	}
}

func (m *SimConnectManager) connect() {
	m.logInfo("[SimConnectManager] Attempting to connect...")
	m.client = client.New("wails")
	err := m.client.Connect()
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if err != nil {
		m.logDebug(fmt.Sprintf("[SimConnectManager] Connection failed: %v", err))
		m.state = Offline
		m.setConnected(false)
		return
	}
	m.logInfo("[SimConnectManager] Connected successfully.")
	m.state = Online
	m.setConnected(true)
	go m.listen()
	go m.monitorSystemState()
}

// monitorSystemState requests system state every second and checks for connection loss
func (m *SimConnectManager) monitorSystemState() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		m.stateMu.Lock()
		if m.state != Online || m.client == nil {
			m.stateMu.Unlock()
			return
		}
		m.stateMu.Unlock()

		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			err := m.client.RequestSystemStateSim(simStateRequestID)
			if err != nil {
				m.logDebug("[SimConnectManager] System state request failed, treating as disconnect.")
				m.disconnect()
				return
			}
			// Wait for response in listen()
		}
	}
}

func (m *SimConnectManager) disconnect() {
	m.logDebug("[SimConnectManager] Disconnecting...")
	if m.client != nil {
		_ = m.client.Disconnect()
	}
	m.stateMu.Lock()
	m.state = Offline
	m.stateMu.Unlock()
	m.setConnected(false)
	m.logDebug("[SimConnectManager] Disconnected.")
}

func (m *SimConnectManager) listen() {
	responseTimeout := 2 * time.Second
	var lastSimStateResponse time.Time
	for message := range m.client.Stream() {
		if message.Error != nil {
			m.logDebug(fmt.Sprintf("SimConnect error: %v", message.Error))
			continue
		}
		if message.IsQuit() {
			m.logDebug("SimConnect quit signal received")
			m.stateMu.Lock()
			m.state = Offline
			m.stateMu.Unlock()
			m.setConnected(false)
			break
		}
		if message.IsOpen() {
			m.logDebug("SimConnect connection established")
			m.stateMu.Lock()
			m.state = Online
			m.stateMu.Unlock()
			m.setConnected(true)
		}
		// Handle SimConnect messages by type (production pattern)
		switch message.MessageType {
		case types.SIMCONNECT_RECV_ID_SYSTEM_STATE:
			if ev, ok := message.Data.(*types.SIMCONNECT_RECV_SYSTEM_STATE); ok {
				if ev.DwRequestID == simStateRequestID {
					m.simState.SetSim(int(ev.DwInteger))
					lastSimStateResponse = time.Now()
				}
			}
			// Add more cases for other message types as needed
		}
		// Check for missed heartbeat
		if !lastSimStateResponse.IsZero() && time.Since(lastSimStateResponse) > responseTimeout {
			m.logDebug("[SimConnectManager] Missed system state response, treating as disconnect.")
			m.disconnect()
			return
		}
	}
}

func (m *SimConnectManager) setConnected(val bool) {
	select {
	case m.statusCh <- val:
	default:
	}
	// Emit Wails event if context is set
	if m.wailsCtx != nil {
		// Use the Wails runtime to emit an event
		// Event name: "global::sim-status", Data: val (bool)
		go func(v bool) {
			// Avoid blocking
			runtime.EventsEmit(m.wailsCtx, "global::sim-status", v)
		}(val)
	}
}

func (m *SimConnectManager) Status() bool {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	return m.state == Online
}

func (m *SimConnectManager) StatusChan() <-chan bool {
	return m.statusCh
}

func (m *SimConnectManager) logInfo(args ...interface{}) {
	if m.logger != nil {
		msg := fmt.Sprint(args...)
		m.logger.Info(msg)
	}
}

func (m *SimConnectManager) logDebug(args ...interface{}) {
	if m.logger != nil {
		msg := fmt.Sprint(args...)
		m.logger.Debug(msg)
	}
}
