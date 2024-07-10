package local

import (
	"fmt"
	"github.com/allape/goview/datasource/driver"
	"os"
	"testing"
)

func TestDriver(t *testing.T) {
	d := Driver{}

	fmt.Println(os.Getwd())

	var folder = "../../../samples"

	driver.TestDriver(t, &d, folder)
}
