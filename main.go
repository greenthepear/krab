package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func getOutputFilepath(inFp, outFp string) string {
	maybeDir, err := os.Stat(outFp)
	isDir := err == nil && maybeDir != nil && maybeDir.IsDir()
	if outFp == "" || isDir {
		kra := regexp.MustCompile(`.kra$`)
		if kra.MatchString(inFp) {
			inFp = kra.ReplaceAllString(inFp, ".png")
		} else {
			inFp = inFp + ".png"
		}
	}

	if isDir {
		return path.Join(outFp, filepath.Base(inFp))
	}

	if outFp != "" {
		return outFp
	}
	return inFp
}

func main() {
	outputFilepathFlag := flag.String("o", "",
		"Output .png file or directory, omitting will automatically set the output in the same directory with a changed extensionl. Will overwrite")

	continueOnError := flag.Bool("c", false,
		"Do not halt on error, keep process rest of files")
	flag.Parse()

	if flag.Args() == nil {
		log.Fatal("No input files provided")
	}

	deal := func(file string, e error) {
		wrappedError := fmt.Errorf("%s: %w", file, e)
		if !*continueOnError {
			log.Fatal(wrappedError)
		} else {
			log.Println(wrappedError)
		}
	}

	for _, arg := range flag.Args() {
		outpath := getOutputFilepath(arg, *outputFilepathFlag)

		zipReader, err := zip.OpenReader(arg)
		if err != nil {
			deal(arg, err)
			continue
		}
		defer zipReader.Close()

		found := false
		for _, file := range zipReader.File {
			if file.Name != "mergedimage.png" {
				continue
			}
			rc, err := file.Open()
			if err != nil {
				deal(arg, err)
				break
			}
			// We got it
			outputFile, err := os.Create(outpath)
			if err != nil {
				deal(arg, err)
				continue
			}
			defer outputFile.Close()
			_, err = io.Copy(outputFile, rc)
			if err != nil {
				deal(arg, err)
				break
			}
			rc.Close()
			found = true
			break
		}
		if !found {
			deal(arg, errors.New("mergedimage.png not found"))
		}
	}
}
