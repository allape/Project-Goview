package assets

import (
	"testing"
)

const width, height = 320, 480

func TestCreate404(t *testing.T) {
	err := CreateImage(width, height, "i_v_404.jpg", "404 Not Found", 20)
	if err != nil {
		t.Error(err)
	}
}

func TestCreate500(t *testing.T) {
	err := CreateImage(width, height, "i_v_500.jpg", "500 Internal Server Error", 20)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNoPreview(t *testing.T) {
	err := CreateImage(width, height, "i_v_no_preview.jpg", "NO PREVIEW", 20)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateCDImage(t *testing.T) {
	err := CreateImage(width, height, "../ui/src/asset/i_v_cd...jpg", "cd ..", 120)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateFolderImage(t *testing.T) {
	err := CreateImage(width, height, "../ui/src/asset/i_v_folder.jpg", "FOLDER", 80)
	if err != nil {
		t.Error(err)
	}
}
