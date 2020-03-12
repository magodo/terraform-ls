package exec

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if v := os.Getenv("TF_LS_MOCK"); v != "" {
		os.Exit(ExecuteMock(v))
		return
	}

	os.Exit(m.Run())
}