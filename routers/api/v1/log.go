package v1

import (
	"archive/zip"
	"fmt"
	"github.com/AzusaChino/ficus/pkg/pool"
	"github.com/AzusaChino/ficus/service/logging_service"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		_ = c.SendStatus(http.StatusInternalServerError)
	}
	tmpLoc := ""
	// save tmp file
	_ = c.SaveFile(file, tmpLoc+file.Filename)

	// add args to pool
	_ = pool.Pool.Submit(func() {
		logging_service.AsyncSend(file)
	})

	return c.JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "ok",
	})
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func _(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			log.Fatalf("err: %v", err)
		}
	}(r)

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fPath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fPath)
		}

		filenames = append(filenames, fPath)

		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		err = outFile.Close()
		if err != nil {
			return nil, err
		}
		err = rc.Close()
		if err != nil {
			return nil, err
		}

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
