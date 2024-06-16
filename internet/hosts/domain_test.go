package hosts

import (
	"testing"
)

func Test_ReadIps(t *testing.T) {
	ReadIps("github.com.", "8.8.8.8:53")
}
