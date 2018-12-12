package main

import (
 "fmt"
 "io/ioutil"
 "log"
 "os"
 "math"
 "image"
 "image/jpeg"
 "image/png"
 "path/filepath"
    "strings"
    "strconv"
    "image/color"
)

const trainDir = "data/train/"
const sampleDir = "data/samples/"
const RGB = 3
const MAXMIN = 2


var UNDERRYPE = [RGB][MAXMIN] float64 {{0.0,9.3}, {0.0,6.45}, {0.0,4.09}}
var RYPE = [RGB][MAXMIN] float64 {{35.0,118.4}, {14.0,33.3}, {2.8,16.0}}
var OVERRYPE = [RGB][MAXMIN] float64 {{72.1,100.8}, {17.7,31.9}, {21.0,28.0}}
var ABOUTTORYPE = [RGB][MAXMIN] float64 {{6.50,47.9}, {4.1,17.5}, {3.1,19.6}}
var ABOUTTOOVERRIPE = [RGB][MAXMIN] float64 {{83.4,113.8}, {17.6,36.6}, {14.7,25.5}}

var ITT int

var classifierTable = map[string](*[RGB][MAXMIN]float64) {
    "UNDERRYPE": &UNDERRYPE,
    "RYPE": &RYPE,
    "OVERRYPE": &OVERRYPE,
    "ABOUTTORYPE": &ABOUTTORYPE,
    "ABOUTTOOVERRYPE": &ABOUTTOOVERRIPE,
}


//Get rypeness checks that the RGB values are within the ranges specified by the RGB rypness detection paper.
func getRypenessRGB(minR,maxR,minG,maxG,minB,maxB float64) string {
    var tags string
    fmt.Printf("minR %f, maxR %f\n minG %f, maxG %f\n minB %f, maxB %f\n\n",minR,maxR,minG,maxG,minB,maxB)
    for tag := range classifierTable {
        current := classifierTable[tag]
        if minR > current[0][0] && 
           maxR < current[0][1] &&
           minG > current[1][0] && 
           maxG < current[1][1] &&
           minB > current[2][0] && 
           maxB < current[2][1] {
            tags += "_" + tag
        }
    }
    return tags
}
            
        

var prob float64 = 0.001

func main () {
    
    //Only a single argument for the number of itterations to perform on the imaging sharpening step.
    args := os.Args[1:]
    var err error
    ITT = 1
    if len(args) > 0 {
        ITT, err = strconv.Atoi(args[0])
        
        if err != nil {
            ITT = 1
        }
    }
    
    fmt.Println("hello fruitfly")
    //Read in training files.
    trainFilenames := getFilenames(trainDir)
    //sampleFilenames := getFilenames(sampleDir)
    RGBClassifierHistograms := make([][16][4]float64,0)
    for i:= range trainFilenames {
        RGBClassifierHistograms = append(RGBClassifierHistograms,histogramRGBPixels(trainDir + trainFilenames[i]))
    }
    RGBMeanHistogram := calculateMean(RGBClassifierHistograms)
    RGBStdHistogram := calculateStd(RGBClassifierHistograms,RGBMeanHistogram)
    printFHistogram(RGBMeanHistogram)

    //Duplicate scan of YCbCr data
    YCbCrClassifierHistograms := make([][16][4]float64,0)
    for i:= range trainFilenames {
        YCbCrClassifierHistograms = append(YCbCrClassifierHistograms,histogramYCbCrPixels(trainDir + trainFilenames[i]))
    }
    YCbCrMeanHistogram := calculateMean(YCbCrClassifierHistograms)
    YCbCrStdHistogram := calculateStd(YCbCrClassifierHistograms,YCbCrMeanHistogram)
    printFHistogram(YCbCrMeanHistogram)

    processAndPlotAgingData(RGBMeanHistogram,RGBStdHistogram)
    return
    //Training Done, Isolate apple and calculate buckets for 
    sampleFilenames := getFilenames(sampleDir)

    for i := range sampleFilenames {
        rypeness, img := processRGBFile(sampleDir + sampleFilenames[i],RGBMeanHistogram,RGBStdHistogram)
        writeOutProcessedImage(sampleFilenames[i],rypeness, img)
    }
    
    for i := range sampleFilenames {
        rypeness, img := processYCbCrFile(sampleDir + sampleFilenames[i],YCbCrMeanHistogram,YCbCrStdHistogram)
        writeOutProcessedImage(sampleFilenames[i],rypeness, img)
    }

}


//This ratio is the magic of the whole processes. Experiment with it later to get better results
const RATIO = (1.0/16.0)

func sumMeanRGB(mean [16][4]float64, im image.Image, x, y int) float64{
    r, g, b, _ := im.At(x,y).RGBA()
    return ((mean[r>>12][0] * mean[g>>12][1] * mean[b>>12][2])) - (RATIO*RATIO*RATIO)
}

func scoreImageRGB(score [][]float64, bounds image.Rectangle, mean [16][4]float64, im image.Image) ([][]float64){
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            score[x][y] = sumMeanRGB(mean,im,x,y)
        }
    }
    return score
}

func initScore(bounds image.Rectangle) ([][]float64) {
    score := make([][]float64,bounds.Max.X)
    for i := range score {
        score[i] = make([]float64,bounds.Max.Y)
    }
    return score
}

func pagerankImage(score [][]float64, bounds image.Rectangle) ([][]float64) {
    nextScore := initScore(bounds)
    for i:= 0; i < ITT; i++ {
        fmt.Printf("ITT %d\n",i)
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
                    if x < bounds.Max.X - 1 && y > 0 {
                        inScore += score[x+1][y-1]
                    }
                    if x > 0 {
                        inScore += score[x-1][y]
                    }
                    if x < bounds.Max.X -1 {
                        inScore += score[x+1][y]
                    }
                    if x > 0 && y < bounds.Max.Y-1 {
                        inScore += score[x-1][y+1]
                    }
                    if y < bounds.Max.Y -1 {
                        inScore += score[x][y+1]
                    }
                    if y < bounds.Max.Y - 1 && x < bounds.Max.X - 1 {
                        inScore += score[x+1][y+1]
                    }
                    //fmt.Printf("%0.6f",inScore)
                    nextScore[x][y] = inScore
                    //fmt.Println(inScore)
                }

            }
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
                    score[x][y] = score[x][y] +  (0.85) * (nextScore[x][y] / 9.0)
                    score[x][y] = score[x][y] / float64(bounds.Max.X * bounds.Max.Y)
                 }
             }
    }
    return score

}

func getMaxMinRGB(score [][]float64, bounds image.Rectangle, im *image.Image, newimg *image.NRGBA) (minR, maxR, minB, maxB, minG, maxG int) {
    minR, minB, minG = 256,256,256
    maxR, maxB, maxG = 0,0,0
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := (*im).At(x,y).RGBA()
            if score[x][y] > 0.0 {
                (*newimg).Set(x,y, (*im).At(x,y))
                sr := r >>8
                sb := b >>8
                sg := g >>8
                if int(sr) < minR {
                    minR = int(sr)
                }
                if int(sb) < minB {
                    minB = int(sb)
                }
                if int(sg) < minG {
                    minG = int(sg)
                }
                //MAX
                if int(sr) > maxR {
                    maxR = int(sr)
                }
                if int(sb) > maxB {
                    maxB = int(sb)
                }
                if int(sg) > maxG {
                    maxG = int(sg)
                }
                
            } else {
                (*newimg).Set(x,y, color.RGBA{0,0,0,0})
            }
            
        }
    }
    return
}

func writeOutProcessedImage(filename string, rypeness string, newimg *image.NRGBA) {
    base := filepath.Base(filename)
    dir := filepath.Dir(filename)
    ext := filepath.Ext(filename)
    filenamediff := dir + "/" + "_mod_" +rypeness+ base
    log.Printf("Modified File :%s\n",filenamediff)

    f, err := os.Create(filenamediff)
    if err != nil {
        log.Fatal(err)
    }
    
    switch ext {
    case ".png":
        if err := png.Encode(f, newimg); err != nil {
            f.Close()
            log.Fatal(err)
        }
    case ".jpeg":
        if err := jpeg.Encode(f, newimg,&jpeg.Options{Quality: 100}); err != nil {
            f.Close()
            log.Fatal(err)
        }
    case ".jpg":
        if err := jpeg.Encode(f, newimg,&jpeg.Options{Quality: 100}); err != nil {
            f.Close()
            log.Fatal(err)
        }
    }
    if err := f.Close();err !=nil{
        log.Fatal(err)
    }
}

func processRGBFile(filename string, mean, std [16][4]float64) (string, *image.NRGBA) {
    im := openImage(filename)
    bounds := im.Bounds()
    newimg := image.NewNRGBA(image.Rect(bounds.Min.X,bounds.Min.Y,bounds.Max.X,bounds.Max.Y))
    score := initScore(bounds)
    score = scoreImageRGB(score,bounds,mean,im)
    score = pagerankImage(score ,bounds)
    minR, maxR, minB, maxB, minG, maxG := getMaxMinRGB(score, bounds, &im, newimg)
    //TODO seperate rypeness from this function it should be its own invocation
    rypeness := getRypenessRGB(float64(minR),float64(maxR),float64(minG),float64(maxG),float64(minB),float64(maxB))
    return rypeness, newimg
}

func processYCbCrFile(filename string, mean, std [16][4]float64) (string, *image.NRGBA) {
    im := openImage(filename)
    bounds := im.Bounds()
    newimg := image.NewNRGBA(image.Rect(bounds.Min.X,bounds.Min.Y,bounds.Max.X,bounds.Max.Y))
    score := initScore(bounds)
    score = scoreImageRGB(score,bounds,mean,im)
    score = pagerankImage(score ,bounds)
    minR, maxR, minB, maxB, minG, maxG := getMaxMinRGB(score, bounds, &im, newimg)
    //TODO seperate rypeness from this function it should be its own invocation
    rypeness := getRypenessRGB(float64(minR),float64(maxR),float64(minG),float64(maxG),float64(minB),float64(maxB))
    return rypeness, newimg
}

func calculateStd(classifiers [][16][4]float64, mean[16][4]float64) [16][4]float64{
    var std[16][4]float64
    for i := range classifiers {
        for j := range classifiers[i] {
            for k := range classifiers[i][j] {
                d := mean[j][k] - classifiers[i][j][k]
                std[j][k] += d*d
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
    var mean[16][4]float64
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

func processAndPlotAgingData(mean, std [16][4]float64) {
    names := GetAndSortSampleNames("data/Aging_Study_1")
    os.Mkdir("output",os.ModeDir | os.ModePerm)
    for i, applename := range names {
        os.Mkdir(fmt.Sprintf("output/%d",i),os.ModeDir | os.ModePerm)
        for _, photo := range applename {
            rypeness, img := processRGBFile(photo, mean, std)
            path := strings.Split(photo,"/")
            date := path[len(path)-2]
            writeOutProcessedImage(fmt.Sprintf("output/%d/%s.jpg",i,date),rypeness,img)
        }
    }

}

func GetAndSortSampleNames(dir string) [][]string{
    fmt.Println("Geting Sample Names")
    dirs, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }
    sortedNames := make([][]string,0)
    for i, d := range dirs {
        //exceptions to the rule, only grab data samples
        if d.Name() == "read_train" || d.Name() == "name.bash" {
            continue
        }
            
        names := getFilenames(dir + "/" + d.Name())
        if i == 0 {
            for _ = range names {
                sortedNames = append(sortedNames,make([]string,0))
            }
        }
        for i := range names {
            sortedNames[i] = append(sortedNames[i],dir + "/" + d.Name() + "/" + names[i])
        }
    }
    fmt.Println(sortedNames)
    return sortedNames
}
    
func getFilenames(dir string) []string {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }
    names := make([]string,0)
    for _, f := range files {
        names = append(names,f.Name())
    }
    return names
}

func openImage(filename string) image.Image {
    log.Printf("Reading pixels of %s\n",filename);
    reader, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()
    m, _, err := image.Decode(reader)
    if err != nil {
        log.Fatal(err)
    }
    return m
}


//histogram pixels takes in a file and builds histogram for each pixel color
//the buckets are sized at 16, there are individual buckets for each color.
func histogramRGBPixels(filename string) [16][4]float64{
    m := openImage(filename)
    bounds := m.Bounds()
    
    var histogram[16][4]int
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := m.At(x,y).RGBA()
            histogram[r>>12][0]++
            histogram[g>>12][1]++
            histogram[b>>12][2]++
            histogram[a>>12][3]++
        }
    }
    //printHistogram(histogram)

    //normalize
    var histogramf[16][4]float64
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
func histogramYCbCrPixels(filename string) [16][4]float64{
    m := openImage(filename)
    bounds := m.Bounds()
    
    var histogram[16][4]int
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := m.At(x,y).RGBA()
            y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))
            //fmt.Printf("y:%d cb:%d cr:%d\n",y,cb,cr)
            histogram[y>>4][0]++
            histogram[cb>>4][1]++
            histogram[cr>>4][2]++
            //histogram[a>>12][3]++
        }
    }

    //normalize
    var histogramf[16][4]float64
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



