package plopper

var ColorsGreen = []int{16, 22, 28, 34, 40, 46, 15}
var ColorsGrayscale = make([]int, 0, 24)

func init() {
	for i := 23; i >= 0; i-- {
		ColorsGrayscale = append(ColorsGrayscale, 255-i)
	}
}
