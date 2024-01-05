package main

import (
	"bufio"
	"fmt"
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
	log.SetFlags(0)
	if strings.Contains(os.Args[1], "help") {
		fmt.Println("Pass in a source directory and then a list of directories into which you want to sort all images from source")
		fmt.Println("Example: imagesorter ~/Pictures cool_stuff mediocre_stuff can_be_deleted_safely")
		os.Exit(0)
	}
	isSixel := false
	var source string
	sortingDirectories := make([]string, 0, 2)
	for _, arg := range os.Args[1:] {
		switch {
		case arg == "--sixel":
			isSixel = true
		case source == "":
			source = arg
		default:
			sortingDirectories = append(sortingDirectories, arg)
		}
	}
	if len(sortingDirectories) < 2 {
		log.Fatal("Please provide a source directory and at least 2 target directories")
	}

	entries := readFileEntries(source)
	warnings := createSortingDirectories(sortingDirectories)
	question := buildQuestion(sortingDirectories)

	loopOverFiles(source, question, entries, sortingDirectories, warnings, isSixel)
}

func readFileEntries(path string) []os.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("failed to read from provided source directory %s\n%v", path, err)
	}
	return entries
}

func createSortingDirectories(sortingDirectories []string) []string {
	warnings := make([]string, 0, 1)
	for _, sortingDirectory := range sortingDirectories {
		err := os.Mkdir(sortingDirectory, 0750)
		if os.IsExist(err) {
			warning := fmt.Sprintf("warning: directory %s already exists, continuing without an error", sortingDirectory)
			warnings = append(warnings, warning)
			continue
		}
		if err != nil {
			log.Fatalf("failed to create %s\n%v", sortingDirectory, err)
		}
	}
	return warnings
}

func buildQuestion(sortingDirectories []string) string {
	questionBuilder := strings.Builder{}

	questionBuilder.WriteString("where do you want to move the image?\n")
	for i, directory := range sortingDirectories {
		questionBuilder.WriteString(fmt.Sprintf("[%d] %s\n", i+1, directory))
	}

	return questionBuilder.String()
}

func loopOverFiles(source, question string, entries []os.DirEntry, sortingDirectories, warnings []string, isSixel bool) {
	reader := bufio.NewReader(os.Stdin)
	// can't really check, setting true by default
	printer := ImagePrinter{isKittySupported: !isSixel}
	for _, picture := range entries {
		if picture.IsDir() {
			continue
		}
		pictureName := picture.Name()
		if !strings.HasSuffix(pictureName, ".png") &&
			!strings.HasSuffix(pictureName, ".jpeg") &&
			!strings.HasSuffix(pictureName, ".jpg") {
			continue
		}
		screen.Clear()
		screen.MoveTopLeft()

		f := openImageOrFail(source, pictureName)
		img := decodeImageOrFail(f)
		printer.PrintImageOrFail(img)
		_ = f.Close()
		fmt.Printf("%s\n\n", pictureName)
		if len(warnings) != 0 {
			for _, warning := range warnings {
				fmt.Println(warning)
			}
			fmt.Printf("Note that in case of a name conflict, the file in the target directory will be overwritten\n\n")
			warnings = warnings[:0]
		}
		number := checkUserResponse(question, len(sortingDirectories), reader)
		moveFileOrFail(source, sortingDirectories[number-1], pictureName)
	}
}

func openImageOrFail(source, pictureName string) *os.File {
	f, err := os.Open(path.Join(source, pictureName))
	if err != nil {
		log.Fatalf("failed to open %s\n%v", pictureName, err)
	}

	return f
}

func decodeImageOrFail(f *os.File) image.Image {
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode %s\n%v", f.Name(), err)
	}

	return img
}

func checkUserResponse(question string, numberOfOptions int, reader *bufio.Reader) int {
	for {
		fmt.Println(question)
		input, err := reader.ReadString('\n')
		if err != nil {
			// Probably just Ctrl+D
			os.Exit(0)
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
		log.Fatalf("failed to move %s\n%v", pictureName, err)
	}
}
