package rander

import (
	"github.com/linnv/logx"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/png"
	"io"
)

type Background struct {
	Src    *image.Paletted // 背景图像
	Stride int             // 每个格子的长度
}

func InitBackground(stride, widht, hight int, sub image.Image) *Background {
	if stride%2 == 0 { // stride转为奇数
		stride = stride + 1
	}
	src := image.NewPaletted(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: widht * stride, Y: hight * stride},
	}, palette.Plan9)

	// 合成背景图
	for i := 0; i < hight; i++ {
		for j := 0; j < widht; j++ {
			min := image.Point{X: j * stride, Y: i * stride}
			max := min.Add(image.Point{X: stride, Y: stride})
			rectangle := image.Rectangle{Min: min, Max: max}
			draw.Draw(src, rectangle, sub, image.Point{}, draw.Src)
		}
	}
	return &Background{
		Src:    src,
		Stride: stride,
	}
}

// 白底空心盒子
func GenerateHollowBox(stride int, c color.Color) image.Image {
	r := image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: stride, Y: stride},
	}
	box := image.NewRGBA(r)
	// 设置边框颜色
	for i := 0; i < stride; i++ {
		for j := 0; j < stride; j++ {
			if i == 0 || i == stride-1 || j == 0 || j == stride-1 {
				box.Set(j, i, c)
			} else {
				box.Set(j, i, color.White)
			}
		}
	}
	return box
}

// 实心盒子
func GenerateSolidBox(stride int, c color.Color) image.Image {
	r := image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: stride, Y: stride},
	}
	box := image.NewPaletted(r, color.Palette{c})
	return box
}

// Rand 渲染一副图片
func (b *Background) Rand(w io.Writer) error {
	return png.Encode(w, b.Src)
}

func (b *Background) Next(p image.Point, sub image.Image) {
	min := image.Point{X: p.X * b.Stride, Y: p.Y * b.Stride}
	max := min.Add(image.Point{X: b.Stride, Y: b.Stride})
	rectangle := image.Rectangle{Min: min, Max: max}
	draw.Draw(b.Src, rectangle, sub, image.Point{}, draw.Src)
}

func (b *Background) Move(points []image.Point, c color.RGBA) {
	lastPoint := image.Point{}
	up := image.Rectangle{
		Min: image.Point{0, -b.Stride},
		Max: image.Point{1, 1},
	}
	down := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{1, b.Stride + 1},
	}
	left := image.Rectangle{
		Min: image.Point{-b.Stride, 0},
		Max: image.Point{1, 1},
	}
	right := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{b.Stride + 1, 1},
	}
	center := b.Stride/2 + 1
	for i, p := range points {
		if i == 0 {
			lastPoint = image.Point{p.X*b.Stride + center - 1, p.Y*b.Stride + center - 1}
			b.Src.Set(lastPoint.X, lastPoint.Y, c)
		} else {
			nowPoint := image.Point{p.X*b.Stride + center - 1, p.Y*b.Stride + center - 1}
			xSub := nowPoint.X - lastPoint.X
			ySub := nowPoint.Y - lastPoint.Y
			r := image.Rectangle{}
			if xSub == 0 {
				if ySub > 0 { // down
					r = down.Add(lastPoint)
				} else if ySub < 0 { // up
					r = up.Add(lastPoint)
				}
			} else if ySub == 0 {
				if xSub > 0 { // right
					r = right.Add(lastPoint)
				} else if xSub < 0 { // left
					r = left.Add(lastPoint)
				}
			} else {
				logx.Warnf("move err\n")
			}
			sub := image.NewUniform(c)
			draw.Draw(b.Src, r, sub, image.Point{}, draw.Src)
			lastPoint = nowPoint
		}
	}
}
