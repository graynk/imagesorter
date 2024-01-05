package main

import (
	"github.com/dolmen-go/kittyimg"
	"github.com/mattn/go-sixel"
	"image"
	"log"
	"os"
)

type ImagePrinter struct {
	isKittySupported bool
	sixelEncoder     *sixel.Encoder
}

func (ip *ImagePrinter) PrintImageOrFail(img image.Image) {
	var err error
	if ip.isKittySupported {
		err = ip.printKittyImage(img)
	} else {
		err = ip.printSixelImage(img)
	}

	if err != nil {
		log.Fatalf("failed to display the image, check that you're using a terminal that supports "+
			"Kitty's terminal graphics protocol or Sixel\n%v", err)
	}
}

func (ip *ImagePrinter) printKittyImage(img image.Image) error {
	return kittyimg.Fprintln(os.Stdout, img)
}

func (ip *ImagePrinter) printSixelImage(img image.Image) error {
	if ip.sixelEncoder == nil {
		ip.sixelEncoder = sixel.NewEncoder(os.Stdout)
	}
	return ip.sixelEncoder.Encode(img)
}
