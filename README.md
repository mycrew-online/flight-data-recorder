
# Flight Data Recorder

Logging-enabled Wails/Svelte desktop app for MyCrew.online, featuring SimConnect integration and modern Go logging.

---

## Table of Contents

- [About](#about)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [Building](#building)

---

## About

Flight Data Recorder is a cross-platform desktop application built with [Wails](https://wails.io/) and [Svelte](https://svelte.dev/), designed for MyCrew.online. It leverages SimConnect for flight data and uses a unified logging system based on [go-logz](https://github.com/mrlm-net/go-logz).

## Features

- Modern UI with Svelte and Wails
- SimConnect integration for flight data
- Unified, structured logging (go-logz)
- Hot reload for frontend development
- Production builds for Windows, macOS, Linux

## Installation

Clone the repository and install Go (>=1.23) and Node.js (see Wails requirements).

```sh
git clone https://github.com/mycrew-online/flight-data-recorder.git
cd flight-data-recorder
wails doctor
```

## Usage

### Live Development

Run in live development mode:

```sh
wails dev
```

This starts a Vite dev server for fast frontend reloads. For browser-based Go method access, use the dev server at http://localhost:34115.

### Building

To build a redistributable, production mode package:

```sh
wails build
```

---

## Development

- **Logging:** All logs (app, SimConnect, Wails) use [go-logz](https://github.com/mrlm-net/go-logz) via a Wails-compatible adapter. See `internal/logger/` and `internal/logadapter/`.
- **SimConnect:** Connection management and state monitoring in `pkg/simconnect-manager/`.
- **Frontend:** Svelte app in `website/`.
- **Custom Events:** Extend SimConnect event handling in `manager.go` as needed.

### Contributing

Contributions are welcome! Please follow the style and structure of this README and see [go-logz](https://github.com/mrlm-net/go-logz) for inspiration.

---

## License

Apache-2.0
