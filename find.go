package findImagePosition

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"runtime"
	"sync"
	"sync/atomic"
)

var numGoroutines = runtime.NumCPU()

// FindPosition 返回小图在大图中的位置，若未找到返回false
func FindPosition(bigImg, smallImg image.Image) (image.Point, bool) {
	// 将图像转换为RGBA格式
	rgbaBig := imageToRGBA(bigImg)
	rgbaSmall := imageToRGBA(smallImg)

	// 获取尺寸
	wBig, hBig := rgbaBig.Bounds().Dx(), rgbaBig.Bounds().Dy()
	wSmall, hSmall := rgbaSmall.Bounds().Dx(), rgbaSmall.Bounds().Dy()

	if wSmall > wBig || hSmall > hBig {
		return image.Point{}, false
	}

	// 小图每行的字节数
	smallRowBytes := wSmall * 4
	// 小图总字节数
	//smallTotalBytes := wSmall * hSmall * 4
	//smallPix := rgbaSmall.Pix[:smallTotalBytes]

	// 遍历大图中的每个可能的位置
	for y := 0; y <= hBig-hSmall; y++ {
		for x := 0; x <= wBig-wSmall; x++ {
			match := true
			// 逐行比较
			for sy := 0; sy < hSmall; sy++ {
				// 大图当前行的起始位置
				bigRowStart := (y+sy)*rgbaBig.Stride + x*4
				bigRowEnd := bigRowStart + smallRowBytes
				// 小图当前行的起始位置
				smallRowStart := sy * rgbaSmall.Stride
				smallRowEnd := smallRowStart + smallRowBytes

				// 比较行数据
				if !bytes.Equal(rgbaBig.Pix[bigRowStart:bigRowEnd], rgbaSmall.Pix[smallRowStart:smallRowEnd]) {
					match = false
					break
				}
			}
			if match {
				return image.Point{X: x, Y: y}, true
			}
		}
	}

	return image.Point{}, false
}

// FindAnyPosition 单匹配版本,找到任意一个位置就停止查找，若未找到返回false
func FindAnyPosition(bigImg, smallImg image.Image) (image.Point, bool) {
	rgbaBig := imageToRGBA(bigImg)
	rgbaSmall := imageToRGBA(smallImg)

	wBig, hBig := rgbaBig.Bounds().Dx(), rgbaBig.Bounds().Dy()
	wSmall, hSmall := rgbaSmall.Bounds().Dx(), rgbaSmall.Bounds().Dy()

	if wSmall > wBig || hSmall > hBig {
		return image.Point{}, false
	}

	// 预计算小图四个角的uint32值
	var smallCorners [4]uint32
	precomputeCorners(rgbaSmall, wSmall, hSmall, &smallCorners)

	smallRowBytes := wSmall * 4

	var foundFlag uint32
	resultChan := make(chan image.Point, 1)
	maxY := hBig - hSmall
	if maxY < 0 {
		return image.Point{}, false
	}

	//numGoroutines := runtime.NumCPU()
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
				if atomic.LoadUint32(&foundFlag) != 0 {
					return
				}
				for x := 0; x <= wBig-wSmall; x++ {
					if !checkCorners(rgbaBig, x, y, wSmall, hSmall, smallCorners) {
						continue
					}

					if fullMatch(rgbaBig, rgbaSmall, x, y, wSmall, hSmall, smallRowBytes) {
						select {
						case resultChan <- image.Point{X: x, Y: y}:
							atomic.StoreUint32(&foundFlag, 1)
						default:
						}
						return
					}
				}
			}
		}()
	}

	// 等待结果
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	if pt, ok := <-resultChan; ok {
		return pt, true
	}
	return image.Point{}, false
}

// FindAllPositions 返回所有匹配位置，若未找到返回false
func FindAllPositions(bigImg, smallImg image.Image) ([]image.Point, bool) {
	rgbaBig := imageToRGBA(bigImg)
	rgbaSmall := imageToRGBA(smallImg)

	wBig, hBig := rgbaBig.Bounds().Dx(), rgbaBig.Bounds().Dy()
	wSmall, hSmall := rgbaSmall.Bounds().Dx(), rgbaSmall.Bounds().Dy()

	if wSmall > wBig || hSmall > hBig {
		return nil, false
	}

	// 预计算小图四个角的RGB值（忽略Alpha）
	var smallCorners [4]uint32
	precomputeCornersRGB(rgbaSmall, wSmall, hSmall, &smallCorners)

	maxY := hBig - hSmall
	if maxY < 0 {
		return nil, false
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

	//numGoroutines := runtime.NumCPU()
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

	exist := false
	if len(results) > 0 {
		exist = true
	}
	return results, exist
}
