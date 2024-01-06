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
	"slices"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(0)
	if strings.Contains(os.Args[1], "help") {
		fmt.Println("Pass in a source directory and then a list of directories into which you want to sort all images from source")
		fmt.Println("You can pass --sixel arg to fallback from Kitty graphics protocol to Sixels")
		fmt.Println("Example: imagesorter [--sixel] Pictures cool_stuff mediocre_stuff can_be_deleted_safely")
		fmt.Println("You can also pass --scan to change the behaviour of the tool - it will take a single target directory," +
			"and scan the directories inside it, treating them as new target directories")
		os.Exit(0)
	}
	isSixel := false
	isScan := false
	var source string
	sortingDirectories := make([]string, 0, 2)
	for _, arg := range os.Args[1:] {
		switch {
		case arg == "--sixel":
			isSixel = true
		case arg == "--scan":
			isScan = true
		case source == "":
			source = arg
		default:
			sortingDirectories = append(sortingDirectories, arg)
		}
	}

	entries := readFileEntries(source)
	warnings := make([]string, 0, 1)
	if isScan {
		if len(sortingDirectories) != 1 {
			log.Fatalf("when using --scan, you have to provide exactly one target directory")
		}
		target := sortingDirectories[0]
		sortingDirectories = scanSortingDirectories(target)
	} else {
		warnings = createSortingDirectories(sortingDirectories)
	}
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

func scanSortingDirectories(target string) []string {
	dirEntries, err := os.ReadDir(target)
	if err != nil {
		log.Fatalf("Failed to read target dir %s\n%v", target, err)
	}
	sortingDirectories := make([]string, 0, 1)
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		sortingDirectories = append(sortingDirectories, path.Join(target, dirEntry.Name()))
	}
	return sortingDirectories
}

func createSortingDirectories(sortingDirectories []string) []string {
	warnings := make([]string, 0, 1)
	for _, sortingDirectory := range sortingDirectories {
		err := createNewDir(sortingDirectory)
		if err != nil {
			warnings = append(warnings, err.Error())
		}
	}
	return warnings
}

func buildQuestion(sortingDirectories []string) string {
	questionBuilder := strings.Builder{}

	questionBuilder.WriteString("where do you want to move the image? press enter to skip\n")
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
		number, newDir := checkUserResponse(question, len(sortingDirectories), reader)
		if newDir != "" {
			if i := slices.Index(sortingDirectories, newDir); i != -1 {
				number = i
			} else {
				err := createNewDir(newDir)
				if err != nil {
					warnings = append(warnings, err.Error())
				}
				sortingDirectories = append(sortingDirectories, newDir)
				question = buildQuestion(sortingDirectories)
				number = len(sortingDirectories)
			}
		} else if number == 0 {
			continue
		}
		moveFileOrFail(source, sortingDirectories[number-1], pictureName)
	}

	screen.Clear()
	screen.MoveTopLeft()
	fmt.Println("all done")
}

func createNewDir(newDir string) error {
	err := os.Mkdir(newDir, 0750)
	if os.IsExist(err) {
		return fmt.Errorf("warning: directory %s already exists, continuing without an error", newDir)
	} else if err != nil {
		log.Fatalf("failed to create %s\n%v", newDir, err)
	}

	return nil
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

func checkUserResponse(question string, numberOfOptions int, reader *bufio.Reader) (int, string) {
	for {
		fmt.Println(question)
		input, err := reader.ReadString('\n')
		if err != nil {
			// Probably just Ctrl+D
			os.Exit(0)
		}
		if input == "\n" {
			return 0, ""
		}
		input = strings.TrimSpace(input)
		number, err := strconv.Atoi(input)
		if err != nil {
			return 0, input
		}
		if number <= 0 || number > numberOfOptions {
			fmt.Printf("Please enter a number between 1 and %d\n", numberOfOptions)
			continue
		}
		return number, ""
	}
}

func moveFileOrFail(source, target, pictureName string) {
	err := os.Rename(path.Join(source, pictureName), path.Join(target, pictureName))
	if err != nil {
		log.Fatalf("failed to move %s\n%v", pictureName, err)
	}
}
