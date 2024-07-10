package driver

import (
	"bytes"
	"testing"
)

func TestDriver(t *testing.T, driver Driver, wd string) {
	name := driver.PathJoin(wd, "1.txt")
	files, err := driver.List(wd)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatalf("Expected to have at least one file in the folder %s", wd)
	}

	file, err := driver.Status(name)
	if err != nil {
		t.Fatal(err)
	}

	if file.Name != "1.txt" {
		t.Fatalf("Expected file name to be 1.txt, got %s", file.Name)
	}
	if file.Size != 4 {
		t.Fatalf("Expected file size to be 4, got %d", file.Size)
	}
	if file.IsDir {
		t.Fatalf("Expected file to be a file, got a directory")
	}

	bufferWriter := &bytes.Buffer{}
	err = driver.Concatenate(name, bufferWriter)
	if err != nil {
		t.Fatal(err)
	}
	if string(bufferWriter.Bytes()) != "123\n" {
		t.Fatalf("Expected file content to be 123\\n, got %s", bufferWriter.String())
	}
}
