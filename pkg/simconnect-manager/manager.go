package simconnectmanager

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
	"unsafe"

	logz "github.com/mrlm-net/go-logz/pkg/logger"
	"github.com/mrlm-net/simconnect/pkg/client"
	"github.com/mrlm-net/simconnect/pkg/types"
	"github.com/mycrew-online/flight-data-recorder/internal/logadapter"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// bytesToString converts a null-terminated byte array to a Go string
func bytesToString(b []byte) string {
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

// AirplaneData matches the simvar order and types for SimConnect data definition
type AirplaneData struct {
	Title     [256]byte // string256, units: blank
	Latitude  float64   // radians
	Longitude float64   // radians
	Altitude  float64   // feet
	Heading   float64   // radians
	Airspeed  float64   // knots
	OnGround  float64   // bool as float64 (0/1)
}

// AirplaneState holds the main simvars to be monitored and is extensible for future fields
type AirplaneState struct {
	Title     string  `json:"title"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Heading   float64 `json:"heading"`
	Airspeed  float64 `json:"airspeed"`
	OnGround  bool    `json:"on_ground"`
	// Add more fields as needed for future extension
}

// EnvironmentData matches the simvar order and types for environment data definition
// Using packed struct to match SimConnect's memory layout exactly
type EnvironmentData struct {
	ZuluTime       int32 // seconds since midnight
	LocalTime      int32 // seconds since midnight (local)
	SimTime        int32 // seconds since sim start (int32)
	ZuluDay        int32 // day of month (Zulu)
	ZuluMonth      int32 // month of year (Zulu)
	ZuluYear       int32 // year (Zulu)
	LocalDay       int32 // day of month (Local)
	LocalMonth     int32 // month of year (Local)
	LocalYear      int32 // year (Local)
	ZuluDayOfWeek  int32 // day of week (Zulu)
	LocalDayOfWeek int32 // day of week (Local)
	// Weather variables - no manual padding, let Go handle alignment
	SeaLevelPressure     float64 // SEA LEVEL PRESSURE inHg
	AmbientTemperature   float64 // AMBIENT TEMPERATURE in Celsius
	AmbientWindDirection float64 // AMBIENT WIND DIRECTION in degrees
	AmbientWindVelocity  float64 // AMBIENT WIND VELOCITY in knots
	AmbientVisibility    float64 // AMBIENT VISIBILITY in meters
}

// EnvironmentState holds the main environment vars to be monitored
type EnvironmentState struct {
	ZuluTime       int32 `json:"zulu_time"`
	LocalTime      int32 `json:"local_time"`
	SimTime        int32 `json:"sim_time"`
	ZuluDay        int32 `json:"zulu_day"`
	ZuluMonth      int32 `json:"zulu_month"`
	ZuluYear       int32 `json:"zulu_year"`
	LocalDay       int32 `json:"local_day"`
	LocalMonth     int32 `json:"local_month"`
	LocalYear      int32 `json:"local_year"`
	ZuluDayOfWeek  int32 `json:"zulu_day_of_week"`
	LocalDayOfWeek int32 `json:"local_day_of_week"`
	// Weather variables
	SeaLevelPressure     float64 `json:"sea_level_pressure"`
	AmbientTemperature   float64 `json:"ambient_temperature"`
	AmbientWindDirection float64 `json:"ambient_wind_direction"`
	AmbientWindVelocity  float64 `json:"ambient_wind_velocity"`
	AmbientVisibility    float64 `json:"ambient_visibility"`
	// Add more fields as needed for future extension
}

const simStateRequestID uint32 = 1001

// --- SimulatorState for system state monitoring ---
type SimulatorState struct {
	mu             sync.RWMutex
	Sim            int
	Pause          int
	Crashed        int
	View           int
	AircraftLoaded string
	FlightLoaded   string
	FlightPlan     string
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
	client           *client.Engine
	state            int
	stateMu          sync.Mutex
	stopCh           chan struct{}
	stopped          sync.WaitGroup
	statusCh         chan bool // true=connected, false=disconnected
	logger           *logadapter.LogzWailsAdapter
	simState         SimulatorState
	airplaneState    AirplaneState
	environmentState EnvironmentState
	wailsCtx         context.Context // Wails context for event emission
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
	m.client = client.New("MyCrew.online FDR")
	err := m.client.Connect()
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	if err != nil {
		m.logDebug(fmt.Sprintf("[SimConnectManager] Connection failed: %v", err))
		m.state = Offline
		m.setConnected(false)
		return
	}
	// Register simvar data definition (matches AirplaneData struct)
	defineID := 1
	_ = m.client.AddToDataDefinition(defineID, "TITLE", "", types.SIMCONNECT_DATATYPE_STRING256, 0.0, 0)
	_ = m.client.AddToDataDefinition(defineID, "PLANE LATITUDE", "radians", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 1)
	_ = m.client.AddToDataDefinition(defineID, "PLANE LONGITUDE", "radians", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 2)
	_ = m.client.AddToDataDefinition(defineID, "PLANE ALTITUDE", "feet", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 3)
	_ = m.client.AddToDataDefinition(defineID, "PLANE HEADING DEGREES TRUE", "radians", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 4)
	_ = m.client.AddToDataDefinition(defineID, "AIRSPEED INDICATED", "knots", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 5)
	_ = m.client.AddToDataDefinition(defineID, "SIM ON GROUND", "bool", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 6)
	// Register environment data definition (matches EnvironmentData struct)
	envDefineID := 2
	_ = m.client.AddToDataDefinition(envDefineID, "ZULU TIME", "seconds", types.SIMCONNECT_DATATYPE_INT32, 0.0, 0)
	_ = m.client.AddToDataDefinition(envDefineID, "LOCAL TIME", "seconds", types.SIMCONNECT_DATATYPE_INT32, 0.0, 1)
	_ = m.client.AddToDataDefinition(envDefineID, "SIMULATION TIME", "seconds", types.SIMCONNECT_DATATYPE_INT32, 0.0, 2)
	_ = m.client.AddToDataDefinition(envDefineID, "ZULU DAY OF MONTH", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 3)
	_ = m.client.AddToDataDefinition(envDefineID, "ZULU MONTH OF YEAR", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 4)
	_ = m.client.AddToDataDefinition(envDefineID, "ZULU YEAR", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 5)
	_ = m.client.AddToDataDefinition(envDefineID, "LOCAL DAY OF MONTH", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 6)
	_ = m.client.AddToDataDefinition(envDefineID, "LOCAL MONTH OF YEAR", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 7)
	_ = m.client.AddToDataDefinition(envDefineID, "LOCAL YEAR", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 8)
	_ = m.client.AddToDataDefinition(envDefineID, "ZULU DAY OF WEEK", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 9)
	_ = m.client.AddToDataDefinition(envDefineID, "LOCAL DAY OF WEEK", "number", types.SIMCONNECT_DATATYPE_INT32, 0.0, 10)
	// Weather variables
	_ = m.client.AddToDataDefinition(envDefineID, "SEA LEVEL PRESSURE", "inHg", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 11)
	_ = m.client.AddToDataDefinition(envDefineID, "AMBIENT TEMPERATURE", "celsius", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 12)
	_ = m.client.AddToDataDefinition(envDefineID, "AMBIENT WIND DIRECTION", "degrees", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 13)
	_ = m.client.AddToDataDefinition(envDefineID, "AMBIENT WIND VELOCITY", "knots", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 14)
	_ = m.client.AddToDataDefinition(envDefineID, "AMBIENT VISIBILITY", "meters", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 15)
	// Request data on user aircraft every sim frame
	err = m.client.RequestDataOnSimObject(1, defineID, 0, types.SIMCONNECT_PERIOD_SECOND, types.SIMCONNECT_DATA_REQUEST_FLAG_CHANGED, 0, 0, 0)
	// Request environment data every sim frame
	err2 := m.client.RequestDataOnSimObject(2, envDefineID, 0, types.SIMCONNECT_PERIOD_SECOND, types.SIMCONNECT_DATA_REQUEST_FLAG_CHANGED, 0, 0, 0)
	if err != nil {
		m.logDebug("Failed to request simvar data:", err)
	}
	if err2 != nil {
		m.logDebug("Failed to request environment data:", err2)
	}
	m.logInfo("[SimConnectManager] Connected successfully.")
	m.state = Online
	m.setConnected(true)
	// Subscribe to system events for live updates
	_ = m.client.SubscribeToSystemEvent(100, "Pause")
	_ = m.client.SubscribeToSystemEvent(101, "AircraftLoaded")
	_ = m.client.SubscribeToSystemEvent(102, "FlightLoaded")
	_ = m.client.SubscribeToSystemEvent(103, "Crashed")
	_ = m.client.SubscribeToSystemEvent(107, "Sim")
	_ = m.client.SubscribeToSystemEvent(108, "View")

	// Request initial system state values (one-shot, not heartbeat)
	if err := m.requestInitialSystemStates(); err != nil {
		m.logDebug("Failed to request initial system states:", err)
	}

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
		case types.SIMCONNECT_RECV_ID_EVENT:
			if ev, ok := message.Data.(*types.SIMCONNECT_RECV_EVENT); ok {
				updated := false
				m.simState.mu.Lock()
				switch ev.UEventID {
				case 100: // Pause
					m.simState.Pause = int(ev.DwData)
					updated = true
				case 101: // AircraftLoaded
					// No string data, handled by SYSTEM_STATE
				case 102: // FlightLoaded
					// No string data, handled by SYSTEM_STATE
				case 103: // Crashed
					m.simState.Crashed = int(ev.DwData)
					updated = true
				case 107: // Sim
					m.simState.Sim = int(ev.DwData)
					updated = true
				case 108: // View
					m.simState.View = int(ev.DwData)
					updated = true
				}
				m.simState.mu.Unlock()
				// Emit simulator state to frontend if updated
				if updated {
					m.simState.mu.RLock()
					stateCopy := struct {
						Sim            int    `json:"Sim"`
						Pause          int    `json:"Pause"`
						Crashed        int    `json:"Crashed"`
						View           int    `json:"View"`
						AircraftLoaded string `json:"AircraftLoaded"`
						FlightLoaded   string `json:"FlightLoaded"`
						FlightPlan     string `json:"FlightPlan"`
					}{
						Sim:            m.simState.Sim,
						Pause:          m.simState.Pause,
						Crashed:        m.simState.Crashed,
						View:           m.simState.View,
						AircraftLoaded: m.simState.AircraftLoaded,
						FlightLoaded:   m.simState.FlightLoaded,
						FlightPlan:     m.simState.FlightPlan,
					}
					m.simState.mu.RUnlock()
					m.logInfo("SimulatorState: ", stateCopy)
					if m.wailsCtx != nil {
						go func(state interface{}) {
							runtime.EventsEmit(m.wailsCtx, "simulator::state", state)
						}(stateCopy)
					}
				}
			}
		case types.SIMCONNECT_RECV_ID_SYSTEM_STATE:
			if ev, ok := message.Data.(*types.SIMCONNECT_RECV_SYSTEM_STATE); ok {
				var updated bool
				switch ev.DwRequestID {
				case simStateRequestID:
					m.simState.SetSim(int(ev.DwInteger))
					lastSimStateResponse = time.Now()
					updated = true
				case 101: // AircraftLoaded
					m.simState.mu.Lock()
					m.simState.AircraftLoaded = bytesToString(ev.SzString[:])
					m.simState.mu.Unlock()
					updated = true
				case 102: // FlightLoaded
					m.simState.mu.Lock()
					m.simState.FlightLoaded = bytesToString(ev.SzString[:])
					m.simState.mu.Unlock()
					updated = true
				case 103: // FlightPlan
					m.simState.mu.Lock()
					m.simState.FlightPlan = bytesToString(ev.SzString[:])
					m.simState.mu.Unlock()
					updated = true
				case 104: // Sim (one-shot)
					m.simState.SetSim(int(ev.DwInteger))
					updated = true
				}
				// Emit simulator state to frontend if updated
				if updated && m.wailsCtx != nil {
					// Copy state under lock
					m.simState.mu.RLock()
					stateCopy := struct {
						Sim            int    `json:"Sim"`
						Pause          int    `json:"Pause"`
						Crashed        int    `json:"Crashed"`
						View           int    `json:"View"`
						AircraftLoaded string `json:"AircraftLoaded"`
						FlightLoaded   string `json:"FlightLoaded"`
						FlightPlan     string `json:"FlightPlan"`
					}{
						Sim:            m.simState.Sim,
						Pause:          m.simState.Pause,
						Crashed:        m.simState.Crashed,
						View:           m.simState.View,
						AircraftLoaded: m.simState.AircraftLoaded,
						FlightLoaded:   m.simState.FlightLoaded,
						FlightPlan:     m.simState.FlightPlan,
					}
					m.simState.mu.RUnlock()
					go func(state interface{}) {
						runtime.EventsEmit(m.wailsCtx, "simulator::state", state)
					}(stateCopy)
				}
			}
		case types.SIMCONNECT_RECV_ID_SIMOBJECT_DATA:
			if data, ok := message.Data.(*types.SIMCONNECT_RECV_SIMOBJECT_DATA); ok {
				switch data.DwDefineID {
				case 1:
					airplaneData := (*AirplaneData)(unsafe.Pointer(&data.DwData))
					// Convert [256]byte TITLE field to string (remove null terminators)
					title := string(airplaneData.Title[:])
					for i := range airplaneData.Title {
						if airplaneData.Title[i] == 0 {
							title = string(airplaneData.Title[:i])
							break
						}
					}
					// Assign to AirplaneState
					m.airplaneState.Title = title
					m.airplaneState.Latitude = airplaneData.Latitude * 180.0 / math.Pi
					m.airplaneState.Longitude = airplaneData.Longitude * 180.0 / math.Pi
					m.airplaneState.Altitude = airplaneData.Altitude
					m.airplaneState.Heading = airplaneData.Heading * 180.0 / math.Pi
					m.airplaneState.Airspeed = airplaneData.Airspeed
					m.airplaneState.OnGround = airplaneData.OnGround > 0.5

					m.logInfo("AirplaneState: ", m.airplaneState)
					// Emit airplane state to frontend
					if m.wailsCtx != nil {
						go func(state AirplaneState) {
							runtime.EventsEmit(m.wailsCtx, "airplane::state", state)
						}(m.airplaneState)
					}
				case 2:
					// Parse environment data manually due to struct alignment issues
					// SimConnect packs data tightly, but Go adds padding for alignment
					dataPtr := unsafe.Pointer(&data.DwData)

					// Parse int32 fields (11 fields, 4 bytes each = 44 bytes)
					m.environmentState.ZuluTime = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 0))
					m.environmentState.LocalTime = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 4))
					m.environmentState.SimTime = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 8))
					m.environmentState.ZuluDay = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 12))
					m.environmentState.ZuluMonth = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 16))
					m.environmentState.ZuluYear = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 20))
					m.environmentState.LocalDay = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 24))
					m.environmentState.LocalMonth = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 28))
					m.environmentState.LocalYear = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 32))
					m.environmentState.ZuluDayOfWeek = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 36))
					m.environmentState.LocalDayOfWeek = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 40))

					// Parse float64 fields starting at byte 44 (no padding in SimConnect data)
					m.environmentState.SeaLevelPressure = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 44))
					m.environmentState.AmbientTemperature = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 52))
					m.environmentState.AmbientWindDirection = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 60))
					m.environmentState.AmbientWindVelocity = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 68))
					m.environmentState.AmbientVisibility = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 76))

					m.logInfo("EnvironmentState: ", m.environmentState)

					// Emit environment state to frontend
					if m.wailsCtx != nil {
						go func(state EnvironmentState) {
							runtime.EventsEmit(m.wailsCtx, "environment::state", state)
						}(m.environmentState)
					}
				}
			}
		}
		// Check for missed heartbeat
		if !lastSimStateResponse.IsZero() && time.Since(lastSimStateResponse) > responseTimeout {
			m.logDebug("[SimConnectManager] Missed system state response, treating as disconnect.")
			m.disconnect()
			return
		}
	}
}

// GetAirplaneState returns a copy of the current airplane state
func (m *SimConnectManager) GetAirplaneState() AirplaneState {
	return m.airplaneState
}

// GetEnvironmentState returns a copy of the current environment state
func (m *SimConnectManager) GetEnvironmentState() EnvironmentState {
	return m.environmentState
}

func (m *SimConnectManager) GetSimulatorState() SimulatorState {
	m.simState.mu.RLock()
	defer m.simState.mu.RUnlock()
	return m.simState
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

// requestInitialSystemStates requests AircraftLoaded, FlightLoaded, FlightPlan, Sim (one-shot, not heartbeat)
func (m *SimConnectManager) requestInitialSystemStates() error {
	if m.client == nil {
		return fmt.Errorf("SimConnect client not initialized")
	}
	// Use unique request IDs for each
	if err := m.client.RequestSystemStateAircraftLoaded(101); err != nil {
		return fmt.Errorf("AircraftLoaded request failed: %w", err)
	}
	if err := m.client.RequestSystemStateFlightLoaded(102); err != nil {
		return fmt.Errorf("FlightLoaded request failed: %w", err)
	}
	if err := m.client.RequestSystemStateFlightPlan(103); err != nil {
		return fmt.Errorf("FlightPlan request failed: %w", err)
	}
	if err := m.client.RequestSystemStateSim(104); err != nil {
		return fmt.Errorf("sim request failed: %w", err)
	}
	return nil
}
