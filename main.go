package main

import (
	"archive/zip"
	"bufio"
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
		"Output file or directory, omitting will automatically set the output in the same directory with changed extension. Will overwrite")

	continueOnError := flag.Bool("c", false,
		"Continue (don't halt) on error, keep processing rest of the files")
	stdout := flag.Bool("so", false,
		"Print mergedfile.png into standard output")
	stdinFilename := flag.Bool("sif", false,
		"Get file paths from standard input, needs to be newline seperated")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s: krab [options] [kra file paths]\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `
Examples:
	krab green.kra
	krab -o Graphics/green.png green.kra
	krab -o Graphics/ green.kra
	krab -o Graphics/ green.kra jade.kra mint.kra
	krab -o Graphics/ *.kra
`)
	}
	flag.Parse()
	args := flag.Args()

	if *stdinFilename {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal("reading standard input:", err)
		}
	}

	if len(args) == 0 {
		log.Fatal("no input file(s) provided")
	}

	if len(args) == 1 {
		// In case unexpanded glob gets passed
		// If you name a file with wildcards that's on you
		var err error
		fn := args[0]
		args, err = filepath.Glob(fn)
		if err != nil {
			log.Fatal(err)
		}
		if len(args) == 0 {
			log.Fatal("file does not exist: ", fn)
		}
	}

	deal := func(file string, e error) {
		wrappedError := fmt.Errorf("%s: %w", file, e)
		if !*continueOnError {
			log.Fatal(wrappedError)
		} else {
			log.Println(wrappedError)
		}
	}

	for _, arg := range args {
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
			outputFile := os.Stdout
			if !*stdout {
				outputFile, err = os.Create(outpath)
				if err != nil {
					deal(arg, err)
					continue
				}
				defer outputFile.Close()
			}
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
