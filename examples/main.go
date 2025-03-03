package main

import (
	"fmt"
	fip "github.com/topascend/findImagePosition"
	"image"

	"os"
	"time"
)

func main() {

	path := "small"
	// 测试性能
	for i := 0; i < 50; i++ {
		path = fmt.Sprintf("examples\\small%d.png", i%5)
		begin := time.Now()
		findImage("examples\\big.png", path)
		fmt.Printf("picture: %v taked total time: : %v \n\n", path, time.Since(begin).String())
	}

}

func findImage(bigPath, smallPath string) {
	// 打开大图
	bigFile, err := os.Open(bigPath)
	if err != nil {
		panic(err)
	}
	defer bigFile.Close()
	bigImg, _, err := image.Decode(bigFile)
	if err != nil {
		panic(err)
	}

	// 打开小图
	smallFile, err := os.Open(smallPath)
	if err != nil {
		panic(err)
	}
	defer smallFile.Close()
	smallImg, _, err := image.Decode(smallFile)
	if err != nil {
		panic(err)
	}

	begin := time.Now()
	pos, found := fip.FindPosition(bigImg, smallImg)
	fmt.Printf("findPosition taked time:  %v \n", time.Since(begin).String())
	if found {
		fmt.Println("Found at:", pos.X, pos.Y)
	} else {
		fmt.Println("Not found")
	}

	begin = time.Now()
	pos, found = fip.FindAnyPosition(bigImg, smallImg)
	fmt.Printf("findAnyPosition taked time:  %v \n", time.Since(begin).String())
	if found {
		fmt.Println("Found at:", pos.X, pos.Y)
	} else {
		fmt.Println("Not found")
	}

	begin = time.Now()
	pos2, found := fip.FindAllPositions(bigImg, smallImg)
	fmt.Printf("findAllPositions taked time:  %v \n", time.Since(begin).String())
	fmt.Printf("%+v\n", pos2)
	if found {
		fmt.Println("Found at:", pos2)
	} else {
		fmt.Println("Not found")
	}
}
