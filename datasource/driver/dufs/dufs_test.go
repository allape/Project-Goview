package dufs

import (
	"github.com/allape/goview/datasource/driver"
	"testing"
)

// dufs -p 8000 samples
// `samples` folder is in the root of the project
func TestDriver(t *testing.T) {
	d := Driver{}

	err := d.Setup(nil)
	if err != nil {
		t.Fatal(err)
	}

	var folder = "http://127.0.0.1:8000/"

	driver.TestDriver(t, &d, folder)
}
