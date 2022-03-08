package main

import (
	"fmt"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
)

func GenerateMaps(stride int, probability float64) [][]int {
	result := make([][]int, stride)
	for i := 0; i < stride; i++ {
		line := make([]int, stride)
		for j := 0; j < stride; j++ {
			f := rand.Float64()
			if f > probability {
				line[j] = 3
			}
		}
		result[i] = line
	}
	return result
}

func TestGenerateMaps(t *testing.T) {
	res := GenerateMaps(150, 0.5)
	randerInit(res)
	f, _ := os.Create("./background.png")
	png.Encode(f, background.Src)
	for _, line := range res {
		l := make([]string, len(line))
		for i, item := range line {
			l[i] = strconv.Itoa(item)
		}
		fmt.Printf("{%s},\n", strings.Join(l, ","))
	}
}

func TestSqrt(t *testing.T) {
	fmt.Println(int(math.Hypot(float64(1),float64(1))*10))
}