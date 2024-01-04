package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dolmen-go/kittyimg"
	"github.com/inancgumus/screen"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	if strings.Contains(os.Args[1], "help") {
		fmt.Println("Pass in a source directory and then a list of directories into which you want to sort all images from source")
		fmt.Println("Example: imagesorter ~/Pictures cool_stuff mediocre_stuff can_be_deleted_safely")
		os.Exit(0)
	}
	if len(os.Args) < 3 {
		log.Fatal("Please provide a source directory and at least 2 target directories")
	}

	source := os.Args[1]
	entries := readFileEntries(source)

	sortingDirectories := os.Args[2:]
	createSortingDirectories(sortingDirectories)

	question := buildQuestion(sortingDirectories)

	loopOverFiles(source, question, entries, sortingDirectories)
}

func readFileEntries(path string) []os.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		combinedErr := errors.Join(fmt.Errorf("failed to read from provided source directory %s", path), err)
		panic(combinedErr)
	}
	return entries
}

func createSortingDirectories(sortingDirectories []string) {
	for _, sortingDirectory := range sortingDirectories {
		err := os.MkdirAll(sortingDirectory, 0750)
		if os.IsExist(err) {
			fmt.Printf("warning: directory %s already exists, continuning without an error\n", sortingDirectory)
			continue
		}
		if err != nil {
			combinedErr := errors.Join(fmt.Errorf("failed to create %s", sortingDirectory), err)
			panic(combinedErr)
		}
	}
}

func buildQuestion(sortingDirectories []string) string {
	questionBuilder := strings.Builder{}

	questionBuilder.WriteString("where do you want to move the image?\n")
	for i, directory := range sortingDirectories {
		questionBuilder.WriteString(fmt.Sprintf("[%d] %s\n", i+1, directory))
	}

	return questionBuilder.String()
}

func loopOverFiles(source, question string, entries []os.DirEntry, sortingDirectories []string) {
	reader := bufio.NewReader(os.Stdin)
	for _, picture := range entries {
		if picture.IsDir() {
			continue
		}
		pictureName := picture.Name()
		if !strings.HasSuffix(pictureName, ".png") &&
			!strings.HasSuffix(pictureName, ".jpeg") &&
			!strings.HasSuffix(pictureName, ".jpg") {
			fmt.Printf("skipping %s\n", pictureName)
			continue
		}
		screen.Clear()
		screen.MoveTopLeft()

		f := openImageOrFail(source, pictureName)
		img := decodeImageOrFail(f)
		printImageOrFail(img)
		_ = f.Close()
		fmt.Printf("%s\n\n", pictureName)
		number := checkUserResponse(question, len(sortingDirectories), reader)
		moveFileOrFail(source, sortingDirectories[number-1], pictureName)
	}
}

func openImageOrFail(source, pictureName string) *os.File {
	f, err := os.Open(path.Join(source, pictureName))
	if err != nil {
		combinedErr := errors.Join(fmt.Errorf("failed to open %s", pictureName), err)
		panic(combinedErr)
	}

	return f
}

func decodeImageOrFail(f *os.File) image.Image {
	img, _, err := image.Decode(f)
	if err != nil {
		combinedErr := errors.Join(fmt.Errorf("failed to decode %s", f.Name()), err)
		panic(combinedErr)
	}

	return img
}

func printImageOrFail(img image.Image) {
	err := kittyimg.Fprintln(os.Stdout, img)
	if err != nil {
		combinedErr := errors.Join(fmt.Errorf("failed to display the image, check that you're using a terminal that supports terminal graphics protocol"), err)
		panic(combinedErr)
	}
}

func checkUserResponse(question string, numberOfOptions int, reader *bufio.Reader) int {
	for {
		fmt.Println(question)
		input, err := reader.ReadString('\n')
		if err != nil {
			combinedErr := errors.Join(fmt.Errorf("somehow failed to read from stdin"), err)
			panic(combinedErr)
		}
		number, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Printf("Please enter a number between 1 and %d\n", numberOfOptions)
			continue
		}
		return number
	}
}

func moveFileOrFail(source, target, pictureName string) {
	err := os.Rename(path.Join(source, pictureName), path.Join(target, pictureName))
	if err != nil {
		combinedErr := errors.Join(fmt.Errorf("failed to move %s", pictureName), err)
		panic(combinedErr)
	}
}
