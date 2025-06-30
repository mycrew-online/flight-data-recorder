package internal

import (
	"github.com/mycrew-online/flight-data-recorder/internal/engine"
)

func bootstrap() *engine.Engine {
	app := engine.New()

	return app
}
