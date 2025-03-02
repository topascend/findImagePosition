package findImagePosition

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"runtime"
	"sync"
)

// ###################### 公共工具函数 ######################
func precomputeCorners(img *image.RGBA, w, h int, out *[4]uint32) {
	// 左上角
	if len(img.Pix) >= 4 {
		out[0] = binary.LittleEndian.Uint32(img.Pix[0:4])
	}

	// 右上角
	if rightTop := (w - 1) * 4; len(img.Pix) >= rightTop+4 {
		out[1] = binary.LittleEndian.Uint32(img.Pix[rightTop : rightTop+4])
	}

	// 左下角
	if leftBottom := (h - 1) * img.Stride; len(img.Pix) >= leftBottom+4 {
		out[2] = binary.LittleEndian.Uint32(img.Pix[leftBottom : leftBottom+4])
	}

	// 右下角
	if rightBottom := (h-1)*img.Stride + (w-1)*4; len(img.Pix) >= rightBottom+4 {
		out[3] = binary.LittleEndian.Uint32(img.Pix[rightBottom : rightBottom+4])
	}
}

func checkCorners(big *image.RGBA, x, y, wSmall, hSmall int, corners [4]uint32) bool {
	// 左上角
	offset := y*big.Stride + x*4
	if offset+4 > len(big.Pix) || binary.LittleEndian.Uint32(big.Pix[offset:offset+4]) != corners[0] {
		return false
	}

	// 右上角
	offset = y*big.Stride + (x+wSmall-1)*4
	if offset+4 > len(big.Pix) || binary.LittleEndian.Uint32(big.Pix[offset:offset+4]) != corners[1] {
		return false
	}

	// 左下角
	offset = (y+hSmall-1)*big.Stride + x*4
	if offset+4 > len(big.Pix) || binary.LittleEndian.Uint32(big.Pix[offset:offset+4]) != corners[2] {
		return false
	}

	// 右下角
	offset = (y+hSmall-1)*big.Stride + (x+wSmall-1)*4
	if offset+4 > len(big.Pix) || binary.LittleEndian.Uint32(big.Pix[offset:offset+4]) != corners[3] {
		return false
	}

	return true
}

func fullMatch(big, small *image.RGBA, x, y, wSmall, hSmall, smallRowBytes int) bool {
	for sy := 0; sy < hSmall; sy++ {
		bigStart := (y+sy)*big.Stride + x*4
		smallStart := sy * small.Stride

		if !bytes.Equal(
			big.Pix[bigStart:bigStart+smallRowBytes],
			small.Pix[smallStart:smallStart+smallRowBytes],
		) {
			return false
		}
	}
	return true
}

func imageToRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

// findAllPositions 返回所有匹配位置
func findAllPositions(bigImg, smallImg image.Image) []image.Point {
	rgbaBig := imageToRGBA(bigImg)
	rgbaSmall := imageToRGBA(smallImg)

	wBig, hBig := rgbaBig.Bounds().Dx(), rgbaBig.Bounds().Dy()
	wSmall, hSmall := rgbaSmall.Bounds().Dx(), rgbaSmall.Bounds().Dy()

	if wSmall > wBig || hSmall > hBig {
		return nil
	}

	// 预计算小图四个角的RGB值（忽略Alpha）
	var smallCorners [4]uint32
	precomputeCornersRGB(rgbaSmall, wSmall, hSmall, &smallCorners)

	maxY := hBig - hSmall
	if maxY < 0 {
		return nil
	}

	resultChan := make(chan image.Point) // 无缓冲通道
	var results []image.Point
	var mu sync.Mutex
	done := make(chan struct{})

	// 结果收集协程
	go func() {
		for pt := range resultChan {
			mu.Lock()
			results = append(results, pt)
			mu.Unlock()
		}
		close(done)
	}()

	numGoroutines := runtime.NumCPU()
	var wg sync.WaitGroup
	taskChan := make(chan int, maxY+1)

	// 生成任务队列
	for y := 0; y <= maxY; y++ {
		taskChan <- y
	}
	close(taskChan)

	// 启动工作协程
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for y := range taskChan {
				for x := 0; x <= wBig-wSmall; x++ {
					// 检查四个角（RGB比较）
					if !checkCornersRGB(rgbaBig, x, y, wSmall, hSmall, smallCorners) {
						continue
					}

					// 完整行比较
					if fullMatchStride(rgbaBig, rgbaSmall, x, y, hSmall) {
						resultChan <- image.Point{X: x, Y: y}
					}
				}
			}
		}()
	}

	wg.Wait()
	close(resultChan)
	<-done // 等待收集完成
	return results
}

// 预计算四个角的RGB值（忽略Alpha）
func precomputeCornersRGB(img *image.RGBA, w, h int, out *[4]uint32) {
	getRGB := func(pix []byte) uint32 {
		if len(pix) < 3 {
			return 0
		}
		return uint32(pix[0])<<16 | uint32(pix[1])<<8 | uint32(pix[2])
	}

	// 左上角
	if len(img.Pix) >= 4 {
		out[0] = getRGB(img.Pix[0:3])
	}

	// 右上角
	if rightTop := (w - 1) * 4; len(img.Pix) >= rightTop+4 {
		out[1] = getRGB(img.Pix[rightTop : rightTop+3])
	}

	// 左下角
	if leftBottom := (h - 1) * img.Stride; len(img.Pix) >= leftBottom+4 {
		out[2] = getRGB(img.Pix[leftBottom : leftBottom+3])
	}

	// 右下角
	if rightBottom := (h-1)*img.Stride + (w-1)*4; len(img.Pix) >= rightBottom+4 {
		out[3] = getRGB(img.Pix[rightBottom : rightBottom+3])
	}
}

// 检查四个角（RGB比较）
func checkCornersRGB(big *image.RGBA, x, y, wSmall, hSmall int, corners [4]uint32) bool {
	getRGB := func(pix []byte) uint32 {
		if len(pix) < 3 {
			return 0
		}
		return uint32(pix[0])<<16 | uint32(pix[1])<<8 | uint32(pix[2])
	}

	// 左上角
	offset := y*big.Stride + x*4
	if offset+4 > len(big.Pix) || getRGB(big.Pix[offset:offset+3]) != corners[0] {
		return false
	}

	// 右上角
	offset = y*big.Stride + (x+wSmall-1)*4
	if offset+4 > len(big.Pix) || getRGB(big.Pix[offset:offset+3]) != corners[1] {
		return false
	}

	// 左下角
	offset = (y+hSmall-1)*big.Stride + x*4
	if offset+4 > len(big.Pix) || getRGB(big.Pix[offset:offset+3]) != corners[2] {
		return false
	}

	// 右下角
	offset = (y+hSmall-1)*big.Stride + (x+wSmall-1)*4
	if offset+4 > len(big.Pix) || getRGB(big.Pix[offset:offset+3]) != corners[3] {
		return false
	}

	return true
}

// 完整行比较（考虑Stride）
func fullMatchStride(big, small *image.RGBA, x, y, hSmall int) bool {
	//wSmall := small.Bounds().Dx()
	smallRowBytes := small.Stride // 使用实际Stride

	for sy := 0; sy < hSmall; sy++ {
		bigStart := (y+sy)*big.Stride + x*4
		smallStart := sy * small.Stride

		// 比较整行数据（包括可能的填充）
		if !bytes.Equal(
			big.Pix[bigStart:bigStart+smallRowBytes],
			small.Pix[smallStart:smallStart+smallRowBytes],
		) {
			return false
		}
	}
	return true
}
