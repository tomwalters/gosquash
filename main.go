/*
This app is designed to serve compressed image assets.

Given a path such as `/uploads/cat.jpg?s=200` the app will check on
disk for a file with the name `cat_s200.jpg`. If this doesn't exist
we check for `cat.jpg` and compress it down to an image of size
200x200 and store this on disk before serving it.

If the `s` paramater is omitted the image will be served at its
natural size.

AllowedFileTypes defines what files the app will work with.
AllowedSizes defines which sizes the app will serve.
UploadsPath defined where to look for files.

TODO: Move config to command line args.
*/

package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

// debug tells the app to output useful dlogs
const debug bool = false

// AllowedFileTypes defines what files the app will work with.
var AllowedFileTypes = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
}

// AllowedSizes defines which sizes the app will serve.
var AllowedSizes = map[int]bool{
	100: true,
	200: true,
	600: true,
}

// UploadsPath defined where to look for files.
const UploadsPath string = "./uploads/"

// Output a message to stdout, but only if debug is true
func dlog(msg string) {
	if debug {
		fmt.Printf("%s\n", msg)
	}
}

/*
Resize a file and save the result to disk. Maintains
aspect ratio when resizing.

source string The existing file path
target string The file path to write the output to
size int The max width of the output

string The output path
err An error which occurred when processing
*/
func resizeFile(source string, target string, size int) (string, error) {
	file, err := os.Open(source)
	if err != nil {
		dlog("Unable to open the file")
		return target, err
	}

	// Decode the file
	file.Seek(0, 0)
	img, _, err := image.Decode(file)

	if err != nil {
		dlog("Unable to decode the file from disk")
		return target, err
	}
	file.Close()

	// Resize it
	m := resize.Resize(uint(size), 0, img, resize.Bicubic)

	// Create a new, empty file
	out, err := os.Create(target)
	if err != nil {
		dlog("Unable to save new file to disk")
		return target, err
	}

	// Make sure the stream is closed once we're done here
	defer out.Close()

	// Save to data to disk
	switch ftype := filepath.Ext(source); strings.ToLower(ftype) {
	case ".jpg", ".jpeg":
		jpeg.Encode(out, m, nil)
	case ".png":
		png.Encode(out, m)
	case ".gif":
		gif.Encode(out, m, nil)
	}

	return target, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := UploadsPath + filepath.Base(r.URL.Path)

		dlog(fmt.Sprintf("Searching for: %s\n", path))

		// Check to see if the file exists
		if _, err := os.Stat(path); err == nil {
			dlog("File exists!")

			// Check the size param
			query := r.URL.Query()
			s, ok := query["s"]

			if !ok {
				// No size param, serve the default file
				http.ServeFile(w, r, path)
			} else {
				// Check that the size is allowed
				// TODO: Handle error
				targetSize, _ := strconv.Atoi(s[0])

				if _, exists := AllowedSizes[targetSize]; exists {
					// Check disk to see if the resized image exists
					ext := filepath.Ext(path)
					resizedpath := strings.Replace(path, ext, fmt.Sprintf("_s%d%s", targetSize, ext), -1)

					if _, err := os.Stat(resizedpath); err == nil {
						dlog("Resized file exists already")

						// Serve the target file
						http.ServeFile(w, r, resizedpath)
					} else {
						// Generate the file
						if _, err := resizeFile(path, resizedpath, targetSize); err != nil {
							// Something went wrong, serve the original path
							http.ServeFile(w, r, path)
						}

						// Serve it
						http.ServeFile(w, r, resizedpath)
					}
				}
			}

		} else {
			dlog("File not found")

			// TODO: Implement default image?
			http.NotFound(w, r)
		}
	})

	http.ListenAndServe(":9990", nil)
}
