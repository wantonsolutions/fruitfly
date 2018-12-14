package fruitfly

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const trainDir = "data/train/"
const sampleDir = "data/samples/"
const RGB = 3
const MAXMIN = 2
const ITT = 1

func Process() {
	//Only a single argument for the number of itterations to perform on the imaging sharpening step.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("hello fruitfly")
	//Read in training files.
	trainFilenames := getFilenames(trainDir)
	//sampleFilenames := getFilenames(sampleDir)
	RGBClassifierHistograms := make([][16][4]float64, 0)
	for i := range trainFilenames {
		RGBClassifierHistograms = append(RGBClassifierHistograms, histogramRGBPixels(trainDir+trainFilenames[i]))
	}
	RGBMeanHistogram := calculateMean(RGBClassifierHistograms)
	RGBStdHistogram := calculateStd(RGBClassifierHistograms, RGBMeanHistogram)
	printFHistogram(RGBMeanHistogram)

	//Duplicate scan of YCbCr data
	YCbCrClassifierHistograms := make([][16][4]float64, 0)
	for i := range trainFilenames {
		YCbCrClassifierHistograms = append(YCbCrClassifierHistograms, histogramYCbCrPixels(trainDir+trainFilenames[i]))
	}
	YCbCrMeanHistogram := calculateMean(YCbCrClassifierHistograms)
	YCbCrStdHistogram := calculateStd(YCbCrClassifierHistograms, YCbCrMeanHistogram)
	printFHistogram(YCbCrMeanHistogram)

	processAndPlotAgingData(RGBMeanHistogram, RGBStdHistogram)
	return
	//Training Done, Isolate apple and calculate buckets for
	sampleFilenames := getFilenames(sampleDir)

	for i := range sampleFilenames {
		rypeness, img, _ := processRGBFile(sampleDir+sampleFilenames[i], RGBMeanHistogram, RGBStdHistogram)
		writeOutProcessedImage(sampleFilenames[i], rypeness, img)
	}

	for i := range sampleFilenames {
		rypeness, img := processYCbCrFile(sampleDir+sampleFilenames[i], YCbCrMeanHistogram, YCbCrStdHistogram)
		writeOutProcessedImage(sampleFilenames[i], rypeness, img)
	}

}

func writeOutProcessedImage(filename string, rypeness string, newimg *image.NRGBA) {
	base := filepath.Base(filename)
	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	filenamediff := dir + "/" + rypeness + base
	log.Printf("Modified File :%s\n", filenamediff)

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
		if err := jpeg.Encode(f, newimg, &jpeg.Options{Quality: 100}); err != nil {
			f.Close()
			log.Fatal(err)
		}
	case ".jpg":
		if err := jpeg.Encode(f, newimg, &jpeg.Options{Quality: 100}); err != nil {
			f.Close()
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func processRGBFile(filename string, mean, std [16][4]float64) (string, *image.NRGBA, ColorStat) {
	im := openImage(filename)
	bounds := im.Bounds()
	newimg := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	score := initScore(bounds)
	score = scoreImageRGB(score, bounds, mean, im)
	score = pagerankImage(score, bounds)
	cs := getMaxMinRGB(score, bounds, &im, newimg)
	//TODO seperate rypeness from this function it should be its own invocation
	rypeness := GetRypenessRGB(cs)
	return rypeness, newimg, cs
}

func processYCbCrFile(filename string, mean, std [16][4]float64) (string, *image.NRGBA) {
	im := openImage(filename)
	bounds := im.Bounds()
	newimg := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	score := initScore(bounds)
	score = scoreImageRGB(score, bounds, mean, im)
	score = pagerankImage(score, bounds)
	cs := getMaxMinRGB(score, bounds, &im, newimg)
	//TODO seperate rypeness from this function it should be its own invocation
	rypeness := GetRypenessRGB(cs)
	return rypeness, newimg
}

func processAndPlotAgingData(mean, std [16][4]float64) {
	names := GetAndSortSampleNames("data/Aging_Study_1")
	os.Mkdir("output", os.ModeDir|os.ModePerm)
	aggstatfile, _ := os.Create(fmt.Sprintf("output/stats.dat"))

	for i, applename := range names {
		iplus := i + 1
		os.Mkdir(fmt.Sprintf("output/%d", iplus), os.ModeDir|os.ModePerm)
		statfile, _ := os.Create(fmt.Sprintf("output/%d/stats.dat", iplus))
		for _, photo := range applename {
			rypeness, img, cs := processRGBFile(photo, mean, std)
			path := strings.Split(photo, "/")
			date := path[len(path)-2]
			writeOutProcessedImage(fmt.Sprintf("output/%d/%s.jpg", iplus, date), rypeness, img)
			statfile.WriteString(fmt.Sprintf("%s\n", cs.String()))
			aggstatfile.WriteString(fmt.Sprintf("%d,%s\n", iplus, cs.String()))
		}
	}

}

func GetAndSortSampleNames(dir string) [][]string {
	fmt.Println("Geting Sample Names")
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Dir %s", dir)
		log.Fatal(err)
	}
	sortedNames := make([][]string, 0)
	for i, d := range dirs {
		//exceptions to the rule, only grab data samples
		if d.Name() == "real_train" || d.Name() == "name.bash" {
			continue
		}

		names := getFilenames(dir + "/" + d.Name())
		if i == 0 {
			for _ = range names {
				sortedNames = append(sortedNames, make([]string, 0))
			}
		}
		for i := range names {
			sortedNames[i] = append(sortedNames[i], dir+"/"+d.Name()+"/"+names[i])
		}
		//now sort by date
	}
	for i := range sortedNames {
		sort.Sort(ByDate(sortedNames[i]))
		fmt.Println(sortedNames[i])
	}
	return sortedNames
}

type ByDate []string

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	iStats, _ := os.Stat(a[i])
	jStats, _ := os.Stat(a[j])
	return iStats.ModTime().Before(jStats.ModTime())
}

func getFilenames(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Dir %s", dir)
		log.Fatal(err)
	}
	names := make([]string, 0)
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names
}

func openImage(filename string) image.Image {
	log.Printf("Reading pixels of %s\n", filename)
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
