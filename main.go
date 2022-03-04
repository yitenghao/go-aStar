package main

import (
	"fmt"
	"github.com/linnv/logx"
	"image"
	"image/color"
	"image/gif"
	"os"
	"test/aStar/rander"
)

type Point struct {
	X, Y int
}

var open, close = []element{}, []element{}
var start = element{"", Point{1, 1}, 0, 0, 0}
var end = element{"", Point{1, 8}, 0, 0, 0}
var direction = []string{"↑", "→", "↓", "←"}

var background *rander.Background
var GIF *gif.GIF

type element struct {
	Direction string // ↑ → ↓ ← 方向
	P         Point  // 坐标
	F         int    // 和值  已走步数+预计还需步数
	G         int    // 已走步数
	H         int    // 预计还需步数
}

func Init() {
	stride := 10 // 地图的每个点对应到图像中的格子的长度
	x := len(maps[0])
	y := len(maps)
	box := rander.GenerateSolidBox(stride, color.Black)
	pathBox := rander.GenerateHollowBox(stride, color.RGBA{147, 161, 161, 255})
	startBox := rander.GenerateHollowBox(stride, color.RGBA{0, 0, 255, 255})
	background = rander.InitBackground(stride, x, y, box)

	for i := 0; i < len(maps); i++ {
		line := maps[i]
		for j := 0; j < len(line); j++ {
			switch line[j] {
			case 0:
				background.Next(image.Point{j, i}, pathBox)
			case 1, 2:
				background.Next(image.Point{j, i}, startBox)
			}
		}
	}
	GIF = &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: -1,
	}
	gifAppendImage(GIF, background.Src)
}
func main() {
	defer logx.Flush()
	Init()
	start.G = 0
	start.H = compuntH(start.P, end.P)
	start.F = start.G + start.H
	nextElement := start
	for {
		var err error
		nextElement, err = next(nextElement)
		if err != nil {
			logx.Warnf("next:%+v,%s\n", nextElement, err.Error())
			break
		}
	}

	closePoints := element2ImagePoint(close)
	closeBox := rander.GenerateSolidBox(background.Stride, color.RGBA{0, 255, 0, 255})
	for _, p := range closePoints {
		background.Next(p, closeBox)
		gifAppendImage(GIF, background.Src)
	}

	bestPath := getBestMovePath(close)
	for _, item := range bestPath {
		fmt.Printf("向%s移动，到达%v，和值：%d，步数：%d\n", item.Direction, item.P, item.F, item.G)
	}

	points := element2ImagePoint(bestPath)
	background.Move(points, color.RGBA{255, 0, 0, 255})
	gifAppendImage(GIF, background.Src)

	g, err := os.Create("./move.gif")
	if err != nil {
		logx.Errorf("Create move.gif err:%s\n", err.Error())
	}
	if err = gif.EncodeAll(g, GIF); err != nil {
		logx.Errorf("gif encodeAll err:%s\n", err.Error())
		return
	}
}

func next(s element) (element, error) {
	// 实际移动
	open = Append(open, s)
	open = remove(open, s)
	close = append(close, s)

	if isEqual(s.P, end.P) {
		return s, fmt.Errorf("到达终点")
	}

	for _, d := range direction {
		nextElement, err := elementMove(s, d)
		if err != nil {
			logx.Warnf("elementMove err:%s,nextElement:%+v\n", err.Error(), nextElement)
			continue
		}
		if exist, _ := sliceHas(close, nextElement); exist {
			continue
		}
		if exist, e := sliceHas(open, nextElement); exist {
			if e.F > nextElement.F {
				updateElement(open, nextElement)
			}
		} else {
			open = Append(open, nextElement)
		}
	}
	return open[len(open)-1], nil
}

func remove(list []element, item element) []element {
	for i := len(list) - 1; i >= 0; i-- {
		e := list[i]
		if isEqual(item.P, e.P) {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func Append(list []element, item element) []element {
	equal := -1
	for i, e := range list {
		if e.F >= item.F {
			equal = i
		}
	}
	if equal == -1 {
		list = append([]element{item}, list...)
		return list
	}
	lift := list[:equal]
	right := append([]element{}, list[equal:]...)
	if len(right) > 0 {
		lift = append(lift, right[0])
		right = right[1:]
	}
	list = append(append(lift, item), right...)
	return list
}
func sliceHas(list []element, item element) (bool, element) {
	for _, e := range list {
		if isEqual(item.P, e.P) {
			return true, e
		}
	}
	return false, element{}
}
func updateElement(list []element, item element) {
	for index, e := range list {
		if isEqual(e.P, item.P) {
			list[index] = item
		}
	}
}
func elementMove(s element, direction string) (element, error) {
	nextElement := element{
		Direction: direction,
		G:         s.G + 1,
	}
	switch direction {
	case "↑":
		nextElement.P = Point{
			X: s.P.X - 1,
			Y: s.P.Y,
		}
	case "↓":
		nextElement.P = Point{
			X: s.P.X + 1,
			Y: s.P.Y,
		}
	case "←":
		nextElement.P = Point{
			X: s.P.X,
			Y: s.P.Y - 1,
		}
	case "→":
		nextElement.P = Point{
			X: s.P.X,
			Y: s.P.Y + 1,
		}
	}
	nextElement.H = compuntH(nextElement.P, end.P)
	nextElement.F = nextElement.G + nextElement.H
	if nextElement.P.X > len(maps[0])-1 || nextElement.P.X < 0 || nextElement.P.Y > len(maps)-1 || nextElement.P.Y < 0 {
		return nextElement, fmt.Errorf("边界")
	}
	value := maps[nextElement.P.X][nextElement.P.Y]
	if value == 3 {
		return nextElement, fmt.Errorf("障碍")
	}
	return nextElement, nil
}

// 对比节点
func isEqual(a, b Point) bool {
	if a.X == b.X && a.Y == b.Y {
		return true
	}
	return false
}

// 计算预计还需步数
func compuntH(x, e Point) int {
	a := x.X - e.X
	b := x.Y - e.Y
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	return a + b
}

// 获取最佳移动路径
func getBestMovePath(list []element) []element {
	allPath := make([][]element, 0)
	lastPathIndex := -1
	for _, ele := range list {
		if lastPathIndex != -1 {
			path := allPath[lastPathIndex]
			lastEle := path[len(path)-1]
			if lastEle.G+1 == ele.G && compuntH(lastEle.P, ele.P) == 1 {
				allPath[lastPathIndex] = append(allPath[lastPathIndex], ele)
				continue
			}
		}
		has := false
		for index, path := range allPath {
			lastEle := path[len(path)-1]
			if lastEle.G+1 == ele.G && compuntH(lastEle.P, ele.P) == 1 {
				allPath[index] = append(allPath[index], ele)
				lastPathIndex = index
				has = true
				continue
			}
		}
		if !has {
			left := getLeft(allPath, ele)
			allPath = append(allPath, append(left, ele))
		}

	}
	return allPath[lastPathIndex]
}
func getLeft(allPath [][]element, s element) []element {
	if len(allPath) == 0 {
		return nil
	}
	ele := Point{}
	switch s.Direction {
	case "↑":
		ele = Point{
			X: s.P.X + 1,
			Y: s.P.Y,
		}
	case "↓":
		ele = Point{
			X: s.P.X - 1,
			Y: s.P.Y,
		}
	case "←":
		ele = Point{
			X: s.P.X,
			Y: s.P.Y + 1,
		}
	case "→":
		ele = Point{
			X: s.P.X,
			Y: s.P.Y - 1,
		}
	}
	for _, path := range allPath {
		for i, item := range path {
			if isEqual(item.P, ele) {
				return append([]element{}, path[:i+1]...)
			}
		}
	}
	return nil
}
func element2ImagePoint(eles []element) []image.Point {
	result := make([]image.Point, 0, len(eles))
	for _, e := range eles {
		result = append(result, image.Point{e.P.Y, e.P.X})
	}
	return result
}

func gifAppendImage(g *gif.GIF, p *image.Paletted) {
	pix := make([]uint8, len(p.Pix))
	copy(pix, p.Pix)
	newPaletted := &image.Paletted{
		Pix:     pix,
		Stride:  p.Stride,
		Rect:    p.Rect,
		Palette: p.Palette,
	}

	g.Image = append(g.Image, newPaletted)
	g.Delay = append(g.Delay, 100)
}
