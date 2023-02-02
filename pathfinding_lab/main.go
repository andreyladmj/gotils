package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type LatLon struct {
	LatRad float64
	LonRad float64
}

func main() {
	path := "/Users/andrii.ladyhin/dev/repos/pytils/data/ships_lat_lons.bin"

	file, _ := os.Open(path)

	defer file.Close()

	var latsLons = make([]LatLon, 134526)
	binary.Read(file, binary.LittleEndian, &latsLons)

	fmt.Println(len(latsLons))
	fmt.Println(latsLons[134525])
}
