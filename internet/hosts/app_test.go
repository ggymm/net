package hosts

import (
	"testing"

	"net/pkg/app"
	"net/pkg/log"
)

func TestApp_Devices(t *testing.T) {
	app.Init()
	log.Init("hosts")

	devs := NewApp().Devices()
	for i, dev := range devs {
		t.Log(i, dev)
	}
}
