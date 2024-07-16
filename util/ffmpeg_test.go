package util

import (
	"fmt"
	"image"
	"testing"
)

func TestFFMpegVideoSampleImage(t *testing.T) {
	// NOTE: put a video file into ../samples which is ignored by Git
	videoFile := "../samples/1.m4v"

	output, err := FFMpegVideoSampleImage(videoFile, "../preview/1.m4v.jpg", 10, image.Point{X: 10, Y: 10})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(output))
}

func TestFFMpegScaleImage(t *testing.T) {
	// NOTE: put an image file into ../samples which is ignored by Git
	imageFile := "../samples/1.jpg"

	output, err := FFMpegScaleImage("../preview/1.jpg.jpg", imageFile, 0.1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(output))
}
