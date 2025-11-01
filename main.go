package main

import (
	"archive/zip"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
)

func getOutputFilepath(inFp, outFp string) string {
	if outFp != "" {
		return outFp
	}

	kra := regexp.MustCompile(`.kra$`)
	if kra.MatchString(inFp) {
		return kra.ReplaceAllString(inFp, ".png")
	}
	return inFp + ".png"
}

func main() {
	inputFilepath := flag.String("i", "",
		"Input .kra file")
	outputFilepathFlag := flag.String("o", "",
		"Output .png file, omitting will automatically set the output in the same directory with a changed extension")
	flag.Parse()

	if *inputFilepath == "" {
		log.Fatal("No input file provided (-i)")
	}
	outpath := getOutputFilepath(*inputFilepath, *outputFilepathFlag)

	zipReader, err := zip.OpenReader(*inputFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()

	outputFile, err := os.Create(outpath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	for _, file := range zipReader.File {
		if file.Name != "mergedimage.png" {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(outputFile, rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		os.Exit(0)
	}

	log.Fatal("mergedimage.png not found")
}
