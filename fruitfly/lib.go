package fruitfly

import (
    "os"
	"fmt"
	"image"
	"image/color"
	"math"
    "log"
    "encoding/gob"
)

var UNDERRYPE = [RGB][MAXMIN]float64{{0.0, 9.3}, {0.0, 6.45}, {0.0, 4.09}}
var RYPE = [RGB][MAXMIN]float64{{35.0, 118.4}, {14.0, 33.3}, {2.8, 16.0}}
var OVERRYPE = [RGB][MAXMIN]float64{{72.1, 100.8}, {17.7, 31.9}, {21.0, 28.0}}
var ABOUTTORYPE = [RGB][MAXMIN]float64{{6.50, 47.9}, {4.1, 17.5}, {3.1, 19.6}}
var ABOUTTOOVERRIPE = [RGB][MAXMIN]float64{{83.4, 113.8}, {17.6, 36.6}, {14.7, 25.5}}

var classifierTable = map[string](*[RGB][MAXMIN]float64){
	"UNDERRYPE":       &UNDERRYPE,
	"RYPE":            &RYPE,
	"OVERRYPE":        &OVERRYPE,
	"ABOUTTORYPE":     &ABOUTTORYPE,
	"ABOUTTOOVERRYPE": &ABOUTTOOVERRIPE,
}

//Get rypeness checks that the RGB values are within the ranges specified by the RGB rypness detection paper.
func GetRypenessRGB(cs ColorStat) string {
	var tags string
	fmt.Printf("%s", cs.String())
	for tag := range classifierTable {
		current := classifierTable[tag]
		if float64(cs.minR) > current[0][0] &&
			float64(cs.maxR) < current[0][1] &&
			float64(cs.minG) > current[1][0] &&
			float64(cs.maxG) < current[1][1] &&
			float64(cs.minB) > current[2][0] &&
			float64(cs.maxB) < current[2][1] {
			tags += "_" + tag
		}
	}
	return tags
}

func calculateStd(classifiers [][16][4]float64, mean [16][4]float64) [16][4]float64 {
	var std [16][4]float64
	for i := range classifiers {
		for j := range classifiers[i] {
			for k := range classifiers[i][j] {
				d := mean[j][k] - classifiers[i][j][k]
				std[j][k] += d * d
			}
		}
	}
	for i := range std {
		for j := range std[i] {
			std[i][j] = math.Sqrt(std[i][j] / float64(len(classifiers)))
		}
	}
	return std
}

func calculateMean(classifiers [][16][4]float64) [16][4]float64 {
	var mean [16][4]float64
	for i := range classifiers {
		for j := range classifiers[i] {
			for k := range classifiers[i][j] {
				mean[j][k] += classifiers[i][j][k]
			}
		}
	}
	for i := range mean {
		for j := range mean[i] {
			mean[i][j] = mean[i][j] / float64(len(classifiers))
		}
	}
	return mean
}

//histogram pixels takes in a file and builds histogram for each pixel color
//the buckets are sized at 16, there are individual buckets for each color.
func histogramRGBPixels(filename string) [16][4]float64 {
	m := openImage(filename)
	bounds := m.Bounds()

	var histogram [16][4]int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			histogram[r>>12][0]++
			histogram[g>>12][1]++
			histogram[b>>12][2]++
			histogram[a>>12][3]++
		}
	}
	//printHistogram(histogram)

	//normalize
	var histogramf [16][4]float64
	pixels := (bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X)
	for bucket := range histogram {
		for color := range histogram[bucket] {
			histogramf[bucket][color] = float64(histogram[bucket][color]) / float64(pixels)
		}
	}
	//printFHistogram(histogramf)
	return histogramf
}

//histogram pixels takes in a file and builds histogram for each pixel color
//the buckets are sized at 16, there are individual buckets for each color.
func histogramYCbCrPixels(filename string) [16][4]float64 {
	m := openImage(filename)
	bounds := m.Bounds()

	var histogram [16][4]int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))
			//fmt.Printf("y:%d cb:%d cr:%d\n",y,cb,cr)
			histogram[y>>4][0]++
			histogram[cb>>4][1]++
			histogram[cr>>4][2]++
			//histogram[a>>12][3]++
		}
	}

	//normalize
	var histogramf [16][4]float64
	pixels := (bounds.Max.Y - bounds.Min.Y) * (bounds.Max.X - bounds.Min.X)
	for bucket := range histogram {
		for color := range histogram[bucket] {
			histogramf[bucket][color] = float64(histogram[bucket][color]) / float64(pixels)
		}
	}
	//printFHistogram(histogramf)
	return histogramf
}

func printHistogram(histogram [16][4]int) {
	fmt.Printf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
	for i, x := range histogram {
		fmt.Printf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
	}
}

func printFHistogram(histogramf [16][4]float64) {
	fmt.Printf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
	for i, x := range histogramf {
		fmt.Printf("0x%04x-0x%04x: %0.3f %0.3f %0.3f %0.3f\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
	}
}

//This ratio is the magic of the whole processes. Experiment with it later to get better results
const RATIO = (1.0 / 16.0)

func sumMeanRGB(mean [16][4]float64, im image.Image, x, y int) float64 {
	r, g, b, _ := im.At(x, y).RGBA()
	return (mean[r>>12][0] * mean[g>>12][1] * mean[b>>12][2]) - (RATIO * RATIO * RATIO)
}

func scoreImageRGB(score [][]float64, bounds image.Rectangle, mean [16][4]float64, im image.Image) [][]float64 {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			score[x][y] = sumMeanRGB(mean, im, x, y)
		}
	}
	return score
}

func initScore(bounds image.Rectangle) [][]float64 {
	score := make([][]float64, bounds.Max.X)
	for i := range score {
		score[i] = make([]float64, bounds.Max.Y)
	}
	return score
}

func pagerankImage(score [][]float64, bounds image.Rectangle) [][]float64 {
	nextScore := initScore(bounds)
	for i := 0; i < ITT; i++ {
		fmt.Printf("ITT %d\n", i)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				var inScore float64
				//Directions
				if x > 0 && y > 0 {
					inScore += score[x-1][y-1]
				}
				if y > 0 {
					inScore += score[x][y-1]
				}
				if x < bounds.Max.X-1 && y > 0 {
					inScore += score[x+1][y-1]
				}
				if x > 0 {
					inScore += score[x-1][y]
				}
				if x < bounds.Max.X-1 {
					inScore += score[x+1][y]
				}
				if x > 0 && y < bounds.Max.Y-1 {
					inScore += score[x-1][y+1]
				}
				if y < bounds.Max.Y-1 {
					inScore += score[x][y+1]
				}
				if y < bounds.Max.Y-1 && x < bounds.Max.X-1 {
					inScore += score[x+1][y+1]
				}
				//fmt.Printf("%0.6f",inScore)
				nextScore[x][y] = inScore
				//fmt.Println(inScore)
			}

		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				score[x][y] = score[x][y] + (0.85)*(nextScore[x][y]/9.0)
				score[x][y] = score[x][y] / float64(bounds.Max.X*bounds.Max.Y)
			}
		}
	}
	return score

}

type ColorStat struct {
	minR int
	maxR int
	avgR float64
	stdR float64
	minB int
	maxB int
	avgB float64
	stdB float64
	minG int
	maxG int
	avgG float64
	stdG float64
}

func NewColorStat() ColorStat {
	return ColorStat{}
}

func (cs ColorStat) String() string {
	return fmt.Sprintf("%d,%d,%f,%f,%d,%d,%f,%f,%d,%d,%f,%f", cs.minR, cs.maxR, cs.avgR, cs.stdR, cs.minB, cs.maxB, cs.avgB, cs.stdB, cs.minG, cs.maxG, cs.avgG, cs.stdG)
}

func rollingAverage(m_k_1, x, k float64) (m_k float64) {
	m_k = m_k_1 + (x-m_k_1)/k
	return
}
func rollingStd(s_k_1, x, m_k_1, m_k float64) (s_k float64) {
	s_k = s_k_1 + (x-m_k_1)*(x-m_k)
	return
}

func getMaxMinRGB(score [][]float64, bounds image.Rectangle, im *image.Image, newimg *image.NRGBA) ColorStat {
	cs := NewColorStat()
	cs.minR, cs.minB, cs.minG = 256, 256, 256
	var i int = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := (*im).At(x, y).RGBA()
			if score[x][y] > 0.0 {
				i++

				(*newimg).Set(x, y, (*im).At(x, y))
				sr := r >> 8
				sb := b >> 8
				sg := g >> 8

				//calculate the rolling mean
				var (
					tmpAvgr, tmpAvgb, tmpAvgg float64
				)
				if i == 0 {
					cs.avgR = float64(sr)
					cs.avgG = float64(sg)
					cs.avgB = float64(sb)
				} else {
					//rolling average
					tmpAvgr = rollingAverage(cs.avgR, float64(sr), float64(i))
					tmpAvgg = rollingAverage(cs.avgG, float64(sg), float64(i))
					tmpAvgb = rollingAverage(cs.avgB, float64(sb), float64(i))
					cs.stdR = rollingStd(cs.stdR, float64(sr), cs.avgR, tmpAvgr)
					cs.stdG = rollingStd(cs.stdG, float64(sg), cs.avgG, tmpAvgg)
					cs.stdB = rollingStd(cs.stdB, float64(sb), cs.avgB, tmpAvgb)
					cs.avgR = tmpAvgr
					cs.avgG = tmpAvgg
					cs.avgB = tmpAvgb
				}
				if int(sr) < cs.minR {
					cs.minR = int(sr)
				}
				if int(sb) < cs.minB {
					cs.minB = int(sb)
				}
				if int(sg) < cs.minG {
					cs.minG = int(sg)
				}
				//MAX
				if int(sr) > cs.maxR {
					cs.maxR = int(sr)
				}
				if int(sb) > cs.maxB {
					cs.maxB = int(sb)
				}
				if int(sg) > cs.maxG {
					cs.maxG = int(sg)
				}

			} else {
				(*newimg).Set(x, y, color.RGBA{0, 0, 0, 0})
			}

		}
	}
	return cs
}

func WriteModel(filename string, histogram [16][4]float64) {
    f, err := os.Create(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    enc := gob.NewEncoder(f)
    err = enc.Encode(histogram)
    if err != nil {
        log.Fatal(err)
    }
    return
}

func ReadModel(filename string) [16][4]float64 {
    f, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    var histogram [16][4]float64
    dec := gob.NewDecoder(f)
    err = dec.Decode(&histogram)
    if err != nil {
        log.Fatal(err)
    }
    return histogram
}
