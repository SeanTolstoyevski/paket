package pengine_test

import (
	"testing"
	"github.com/SeanTolstoyevski/paket/pengine"
)

func TestCreateRandomBytesLenght(t *testing.T) {
	bytes, _ := pengine.CreateRandomBytes(32)
	want := 32
	got := len(bytes)
	if got != want {
		t.Error("CreateRandomBytes(16) want 32.")
	}
}

