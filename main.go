package main

/* TODO
Final:
	* Args:
		* directory
		* bucket
		* object prefix? (derive from directory?)
		* thumbnail size(s)
	* For each image in directory:
		* Check if it already exists in bucket
		* resize to specified size
		* upload original and resized both to bucket:prefix/path/within/directory/to/img[_size].jpg

Later:
	* Support multiple resize options
	* Use goroutines for concurrency
	* default path to ./
	* --force to override existance check
*/
import (
	"flag"
	//"fmt"
	"github.com/BurntSushi/toml"
	"github.com/disintegration/imaging"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var params struct {
	size   int
	source string
	bucket string
	prefix string

	verbose bool

	// configuration for object storage
	Access   string
	Secret   string
	Endpoint string
}

func main() {

	objects, err := objectList()
	if err != nil {
		log.Printf("Could not retrieve object list for %s:%s\n", params.bucket, params.prefix)
		log.Fatal(err)
	}

	err = filepath.Walk(params.source, func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}

		objname := pathToObject(path)
		if objects[objname] && objects[objname+":small"] {
			verbose("  %s already exists\n", path)
			return nil
		}

		// FIXME: use goroutine for each path?
		return processImage(path)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func processImage(path string) error {

	origFile, err := os.Open(path)
	if err != nil {
		log.Printf("  **there was an error opening %s: %s", path, err)
		return err
	}
	defer origFile.Close()

	img, format, err := image.Decode(origFile)
	if err != nil {
		verbose("  skipping %s\n", path)
		debug("  **error was: %s\n", err)
		return nil
	}

	//Fit will preserve aspect ratio but resize image to withing bounding box
	img = imaging.Clone(img) // Get it into the right format, just in case
	imgResized := imaging.Fit(img, params.size, params.size, imaging.BSpline)

	// Save resized image to tempfile
	outfile, err := ioutil.TempFile("", "")
	defer os.Remove(outfile.Name())
	defer outfile.Close()

	err = imaging.Encode(outfile, imgResized, getImagingFormat(format))
	if err != nil {
		log.Printf("  **there was an error saving resized image to %s: %s", outfile.Name(), err)
		return nil
	}

	// Upload the original and resized versions to DreamObjects
	origFile.Seek(0, 0)
	uploadImages(origFile, path)

	outfile.Seek(0, 0)
	uploadImages(outfile, path+":small")

	return nil
}

func init() {
	flag.StringVar(&params.bucket, "bucket", "kmarsh", "Name of the bucket to use")
	flag.StringVar(&params.prefix, "prefix", "site/static/img/", "prefix to put on each image uploaded")
	flag.IntVar(&params.size, "size", 1000, "Size of the small image to generate")
	flag.BoolVar(&params.verbose, "verbose", false, "verbosity")

	if _, err := toml.DecodeFile("/Users/kylem/.imgobjsync", &params); err != nil {
		log.Fatal(err)
	}

	flag.Parse() // Parse flags *after* decoding the config file so flags override config
	params.source = flag.Arg(0)
}
