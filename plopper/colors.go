package plopper

type RGBValue [3]int

func GetColorByName(name string) RGBValue {
	switch name {
	case "white":
		return RGBValue{255, 255, 255}
	case "green":
		return RGBValue{0, 255, 0}
	case "red":
		return RGBValue{255, 0, 0}
	case "blue":
		return RGBValue{0, 0, 255}
	case "yellow":
		return RGBValue{255, 255, 0}
	case "orange":
		return RGBValue{255, 200, 0}
	case "dragonberry":
		return RGBValue{255, 40, 167}
	default:
		return RGBValue{255, 0, 255}
	}
}
