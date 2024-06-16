package hosts

import (
	"testing"
	"time"
)

func Test_ReadIps(t *testing.T) {
	ReadIps("github.com.", "8.8.8.8")

	time.Sleep(10 * time.Second)
}
