package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestBathymetry_2160x1080(t *testing.T) {
	// 2160 1080
	file, err := os.Open("data/bathymetry_2160x1080.bin")

	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		log.Fatalf("file stat error: %v", err)
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	if err != nil {
		log.Fatalf("read buffer error: %v", err)
	}

	land := 0
	sea := 0

	for i := 0; i < len(bytes); i++ {
		if bytes[i] == 0 {
			land++
		} else {
			sea++
		}
	}
	fmt.Println("size:", size)
	fmt.Println("land:", land)
	fmt.Println("sea:", sea)
}

func TestBathymetry(t *testing.T) {
	// 21600 10800
	file, err := os.Open("data/bathymetry.bin")

	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		log.Fatalf("file stat error: %v", err)
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	if err != nil {
		log.Fatalf("read buffer error: %v", err)
	}

	land := 0
	sea := 0

	for i := 0; i < len(bytes); i++ {
		if bytes[i] == 0 {
			land++
		} else {
			sea++
		}
	}
	fmt.Println("size:", size)
	fmt.Println("land:", land)
	fmt.Println("sea:", sea)
}
