package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	NameRow = iota
	RedRow
	BlueRow
	GreenRow
	ClassRow
)

type Claim struct {
	Name  string
	Color Pixel
	Count int
	Class int
}

type Pixel struct {
	R int
	G int
	B int
	A int
}

func main() {
	const mapFile, configFile = "map.png", "claims.csv"
	claimsMap := loadImage(mapFile)
	claimsData := loadCSVData(configFile)
	claimsSizes := addClaimSizes(claimsMap, claimsData)
	fmt.Println(claimsSizes)
}

func loadImage(file string) image.Image {
	imageFile, err := os.Open(file)
	if err != nil {
		log.Fatal("Error opening file:", file, err)
	}
	defer imageFile.Close()

	img, err := png.Decode(imageFile)
	if err != nil {
		log.Fatal("Error decoding image:", file, err)
	}

	return img
}

func loadCSVData(file string) []Claim {
	configFile, err := os.Open(file)
	if err != nil {
		log.Fatal("Error opening file:", file, err)
	}
	defer configFile.Close()

	reader := csv.NewReader(configFile)
	records, err := reader.ReadAll()

	claimList := []Claim{}
	for i, record := range records {
		claim, err := recordToClaim(record)
		if i != 0 && err != nil {
			fmt.Println("Could not parse row " + fmt.Sprint(i) + ", skipping.")
		}
		claimList = append(claimList, claim)
	}
	return claimList
}

func recordToClaim(record []string) (Claim, error) {
	red, err := strconv.Atoi(record[RedRow])
	if err != nil {
		return Claim{}, err
	}
	green, err := strconv.Atoi(record[GreenRow])
	if err != nil {
		return Claim{}, err
	}
	blue, err := strconv.Atoi(record[BlueRow])
	if err != nil {
		return Claim{}, err
	}
	class, err := strconv.Atoi(record[ClassRow])
	if err != nil {
		return Claim{}, err
	}

	color := Pixel{
		R: red,
		G: green,
		B: blue,
		A: 255,
	}
	return Claim{
		Name:  record[NameRow],
		Color: color,
		Count: 0,
		Class: class,
	}, nil
}

func addClaimSizes(claimsMap image.Image, claims []Claim) []Claim {
	colorMap := make(map[Pixel]int)
	for x := 0; x <= claimsMap.Bounds().Dx(); x++ {
		for y := 0; y <= claimsMap.Bounds().Dy(); y++ {
			pixelColor := rgbaToPixel(claimsMap.At(x, y).RGBA())
			colorMap[pixelColor] = colorMap[pixelColor] + 1
			for i := 0; i < len(claims); i++ {
				if areColorsEqual(pixelColor, claims[i].Color) {
					claims[i].Count++
				}
			}
		}
	}
	return claims
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r >> 8), int(g >> 8), int(b >> 8), int(a >> 8)}
}

func areColorsEqual(color1, color2 Pixel) bool {
	const ValidOffset = 50
	return math.Abs(float64(color1.R-color2.R)) <= ValidOffset &&
		math.Abs(float64(color1.G-color2.G)) <= ValidOffset &&
		math.Abs(float64(color1.B-color2.B)) <= ValidOffset &&
		math.Abs(float64(color1.A-color2.A)) <= ValidOffset
}
