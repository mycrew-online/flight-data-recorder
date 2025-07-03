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
// The struct layout must match exactly what SimConnect sends
type AirplaneData struct {
	Title           [256]byte // string256, units: blank
	Latitude        float64   // radians
	Longitude       float64   // radians
	Altitude        float64   // feet
	Heading         float64   // radians (true)
	HeadingMagnetic float64   // radians (magnetic)
	Airspeed        float64   // knots
	Bank            float64   // degrees
	AltAboveGround  float64   // feet
	Pitch           float64   // degrees
	VerticalSpeed   float64   // feet/min
}

// AirplaneState holds the main simvars to be monitored and is extensible for future fields
type AirplaneState struct {
	Title           string  `json:"title"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Altitude        float64 `json:"altitude"`
	Heading         float64 `json:"heading"`
	HeadingMagnetic float64 `json:"heading_magnetic"`
	Airspeed        float64 `json:"airspeed"`
	Bank            float64 `json:"bank"`
	AltAboveGround  float64 `json:"alt_above_ground"`
	Pitch           float64 `json:"pitch"`
	VerticalSpeed   float64 `json:"vertical_speed"`
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
	// New simvars
	TimeZoneOffset  int32 // TIME ZONE OFFSET in seconds
	ZuluSunriseTime int32 // ZULU SUNRISE TIME in seconds since midnight
	ZuluSunsetTime  int32 // ZULU SUNSET TIME in seconds since midnight
	TimeOfDay       int32 // TIME OF DAY enum (0=dawn,1=day,2=dusk,3=night)
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
	// New simvars
	TimeZoneOffset  int32 `json:"time_zone_offset"`
	ZuluSunriseTime int32 `json:"zulu_sunrise_time"`
	ZuluSunsetTime  int32 `json:"zulu_sunset_time"`
	TimeOfDay       int32 `json:"time_of_day"`
}

const simStateRequestID uint32 = 1001

// --- SimulatorState for system state monitoring ---
type SimulatorState struct {
	Sim              int     `json:"sim"`
	Pause            int     `json:"pause"`
	Crashed          int     `json:"crashed"`
	View             int     `json:"view"`
	AircraftLoaded   string  `json:"aircraft_loaded"`
	FlightLoaded     string  `json:"flight_loaded"`
	FlightPlan       string  `json:"flight_plan"`
	SimulationRate   float64 `json:"simulation_rate"`
	Realism          int     `json:"realism"`
	SurfaceCondition int     `json:"surface_condition"`
	SurfaceInfoValid int     `json:"surface_info_valid"`
	SurfaceType      int     `json:"surface_type"`
	OnAnyRunway      int     `json:"on_any_runway"`
	InParkingState   int     `json:"in_parking_state"`
	OnGround         bool    `json:"on_ground"`
}

// No mutex or methods needed, match AirplaneState/EnvironmentState style

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
	_ = m.client.AddToDataDefinition(defineID, "PLANE HEADING DEGREES MAGNETIC", "radians", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 5)
	_ = m.client.AddToDataDefinition(defineID, "AIRSPEED INDICATED", "knots", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 6)
	// Remove SIM ON GROUND from AirplaneData definition, add to SimulatorState definition below
	_ = m.client.AddToDataDefinition(defineID, "PLANE BANK DEGREES", "degrees", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 7)
	_ = m.client.AddToDataDefinition(defineID, "PLANE ALT ABOVE GROUND", "feet", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 8)
	_ = m.client.AddToDataDefinition(defineID, "PLANE PITCH DEGREES", "degrees", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 9)
	_ = m.client.AddToDataDefinition(defineID, "VERTICAL SPEED", "feet per minute", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 10)
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
	// New simvars
	_ = m.client.AddToDataDefinition(envDefineID, "TIME ZONE OFFSET", "seconds", types.SIMCONNECT_DATATYPE_INT32, 0.0, 16)
	_ = m.client.AddToDataDefinition(envDefineID, "ZULU SUNRISE TIME", "seconds", types.SIMCONNECT_DATATYPE_INT32, 0.0, 17)
	_ = m.client.AddToDataDefinition(envDefineID, "ZULU SUNSET TIME", "seconds", types.SIMCONNECT_DATATYPE_INT32, 0.0, 18)
	_ = m.client.AddToDataDefinition(envDefineID, "TIME OF DAY", "enum", types.SIMCONNECT_DATATYPE_INT32, 0.0, 19)
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
	// Register additional simvars for SimulatorState
	_ = m.client.AddToDataDefinition(3, "SIMULATION RATE", "", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 0)
	_ = m.client.AddToDataDefinition(3, "REALISM", "", types.SIMCONNECT_DATATYPE_INT32, 0.0, 1)
	_ = m.client.AddToDataDefinition(3, "SURFACE CONDITION", "", types.SIMCONNECT_DATATYPE_INT32, 0.0, 2)
	_ = m.client.AddToDataDefinition(3, "SURFACE INFO VALID", "", types.SIMCONNECT_DATATYPE_INT32, 0.0, 3)
	_ = m.client.AddToDataDefinition(3, "SURFACE TYPE", "", types.SIMCONNECT_DATATYPE_INT32, 0.0, 4)
	_ = m.client.AddToDataDefinition(3, "ON ANY RUNWAY", "", types.SIMCONNECT_DATATYPE_INT32, 0.0, 5)
	_ = m.client.AddToDataDefinition(3, "PLANE IN PARKING STATE", "", types.SIMCONNECT_DATATYPE_INT32, 0.0, 6)
	_ = m.client.AddToDataDefinition(3, "SIM ON GROUND", "bool", types.SIMCONNECT_DATATYPE_FLOAT64, 0.0, 7)
	// Request additional simvars every second
	_ = m.client.RequestDataOnSimObject(3, 3, 0, types.SIMCONNECT_PERIOD_SECOND, types.SIMCONNECT_DATA_REQUEST_FLAG_CHANGED, 0, 0, 0)

	// Request initial system state values (one-shot, not heartbeat)
	if err := m.requestInitialSystemStates(); err != nil {
		m.logDebug("Failed to request initial system states:", err)
	}

	err = m.client.MapClientEventToSimEvent(90111, "PAUSE_ON")
	if err != nil {
		fmt.Println(fmt.Errorf("failed to map PAUSE_TOGGLE event: %v", err))
	}

	err = m.client.AddClientEventToNotificationGroup(1, 90111)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to add event to notification group: %v", err))
	}

	err = m.client.SetNotificationGroupPriority(1, 1000) // High priority
	if err != nil {
		fmt.Println(fmt.Errorf("failed to set notification group priority: %v", err))
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
				// Emit simulator state to frontend if updated
				if updated {
					m.logInfo("SimulatorState: ", m.simState)
					if m.wailsCtx != nil {
						runtime.EventsEmit(m.wailsCtx, "simulator::state", m.simState)
					}
				}
			}
		case types.SIMCONNECT_RECV_ID_SYSTEM_STATE:
			if ev, ok := message.Data.(*types.SIMCONNECT_RECV_SYSTEM_STATE); ok {
				var updated bool
				switch ev.DwRequestID {
				case simStateRequestID:
					m.simState.Sim = int(ev.DwInteger)
					lastSimStateResponse = time.Now()
					updated = true
				case 101: // AircraftLoaded
					m.simState.AircraftLoaded = bytesToString(ev.SzString[:])
					updated = true
				case 102: // FlightLoaded
					m.simState.FlightLoaded = bytesToString(ev.SzString[:])
					updated = true
				case 103: // FlightPlan
					m.simState.FlightPlan = bytesToString(ev.SzString[:])
					updated = true
				case 104: // Sim (one-shot)
					m.simState.Sim = int(ev.DwInteger)
					updated = true
				}
				// Emit simulator state to frontend if updated
				if updated && m.wailsCtx != nil {
					runtime.EventsEmit(m.wailsCtx, "simulator::state", m.simState)
				}
			}
		case types.SIMCONNECT_RECV_ID_SIMOBJECT_DATA:
			if data, ok := message.Data.(*types.SIMCONNECT_RECV_SIMOBJECT_DATA); ok {
				switch data.DwDefineID {
				case 1:
					// Parse airplane data manually from raw bytes to avoid struct padding issues
					dataPtr := unsafe.Pointer(&data.DwData)

					// Title: 256 bytes at offset 0
					titleBytes := (*[256]byte)(unsafe.Pointer(uintptr(dataPtr) + 0))
					m.airplaneState.Title = bytesToString(titleBytes[:])

					// After 256 bytes for title, float64 fields start
					// SimConnect packs data without Go's struct padding
					m.airplaneState.Latitude = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 256)) * 180.0 / math.Pi
					m.airplaneState.Longitude = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 264)) * 180.0 / math.Pi
					m.airplaneState.Altitude = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 272))
					m.airplaneState.Heading = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 280)) * 180.0 / math.Pi
					m.airplaneState.HeadingMagnetic = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 288)) * 180.0 / math.Pi
					m.airplaneState.Airspeed = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 296))
					m.airplaneState.Bank = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 304))
					m.airplaneState.AltAboveGround = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 312))
					m.airplaneState.Pitch = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 320))

					// Extract and format vertical speed
					rawVerticalSpeed := *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 328))
					// Round very small values to zero for cleaner display
					// Values less than 0.1 ft/min are essentially zero for practical purposes
					if math.Abs(rawVerticalSpeed) < 0.1 {
						m.airplaneState.VerticalSpeed = 0.0
					} else {
						m.airplaneState.VerticalSpeed = math.Round(rawVerticalSpeed*100) / 100 // Round to 2 decimal places
					}

					m.logInfo("AirplaneState: ", m.airplaneState)
					// Emit airplane state to frontend
					if m.wailsCtx != nil {
						runtime.EventsEmit(m.wailsCtx, "airplane::state", m.airplaneState)
					}
				case 2:
					// ...existing code for environmentState...
					dataPtr := unsafe.Pointer(&data.DwData)
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
					m.environmentState.SeaLevelPressure = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 44))
					m.environmentState.AmbientTemperature = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 52))
					m.environmentState.AmbientWindDirection = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 60))
					m.environmentState.AmbientWindVelocity = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 68))
					m.environmentState.AmbientVisibility = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 76))
					m.environmentState.TimeZoneOffset = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 84))
					m.environmentState.ZuluSunriseTime = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 88))
					m.environmentState.ZuluSunsetTime = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 92))
					m.environmentState.TimeOfDay = *(*int32)(unsafe.Pointer(uintptr(dataPtr) + 96))
					m.logInfo("EnvironmentState: ", m.environmentState)
					if m.wailsCtx != nil {

						runtime.EventsEmit(m.wailsCtx, "environment::state", m.environmentState)
					}
				case 3:
					// Parse SimulatorState additional simvars
					dataPtr := unsafe.Pointer(&data.DwData)
					m.simState.SimulationRate = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 0))
					m.simState.Realism = int(*(*int32)(unsafe.Pointer(uintptr(dataPtr) + 8)))
					m.simState.SurfaceCondition = int(*(*int32)(unsafe.Pointer(uintptr(dataPtr) + 12)))
					m.simState.SurfaceInfoValid = int(*(*int32)(unsafe.Pointer(uintptr(dataPtr) + 16)))
					m.simState.SurfaceType = int(*(*int32)(unsafe.Pointer(uintptr(dataPtr) + 20)))
					m.simState.OnAnyRunway = int(*(*int32)(unsafe.Pointer(uintptr(dataPtr) + 24)))
					m.simState.InParkingState = int(*(*int32)(unsafe.Pointer(uintptr(dataPtr) + 28)))
					m.simState.OnGround = *(*float64)(unsafe.Pointer(uintptr(dataPtr) + 32)) > 0.5
					m.logInfo("SimulatorState (extra): ", m.simState)
					// Always emit full state to frontend
					if m.wailsCtx != nil {
						runtime.EventsEmit(m.wailsCtx, "simulator::state", m.simState)
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

func (m *SimConnectManager) TogglePause() {
	p := 1 - m.simState.Pause // Toggle pause state

	err := m.client.TransmitClientEvent(
		int(types.SIMCONNECT_OBJECT_ID_USER), // User aircraft
		90111,
		p, // Parameter (external power source 1)
		1,
	)
	if err != nil {
		fmt.Printf("Failed to toggle external power: %v\n", err)
	}
}
