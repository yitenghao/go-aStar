package rander

import (
	"image"
	"image/color"
	"os"
	"testing"
)

func TestInitBackground(t *testing.T) {
	box := GenerateHollowBox(5,color.RGBA{147,161,161,255})
	black := GenerateHollowBox(5,color.Black)
	b := InitBackground(5, 10, 10, black)
	//solidBox:=GenerateSolidBox(b.Stride,color.RGBA{0, 255, 0, 255})
	b.Next(image.Point{1, 2}, box)
	b.Next(image.Point{2, 2}, box)
	b.Next(image.Point{3, 2}, box)
	b.Move([]image.Point{{1, 1}, {1, 0}}, color.RGBA{255, 0, 0, 255}) // up
	b.Move([]image.Point{{1, 1}, {1, 2}}, color.RGBA{255, 0, 0, 255}) // down
	b.Move([]image.Point{{1, 1}, {0, 1}}, color.RGBA{255, 0, 0, 255}) // left
	b.Move([]image.Point{{1, 1}, {2, 1}}, color.RGBA{255, 0, 0, 255}) // right

	f, _ := os.Create("./background.png")
	err := b.Rand(f)
	if err != nil {
		t.Errorf("InitBackground err:%s\n", err.Error())
	}
}
