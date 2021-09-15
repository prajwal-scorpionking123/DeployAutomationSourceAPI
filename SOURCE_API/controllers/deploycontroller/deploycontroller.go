package deploycontroller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/team_six/SOURCE_API/models"
)

//checking whether file in that entered path is exist or not
// func FileExists(name string) bool {
// 	if _, err := os.Stat(name); err != nil {
// 		if os.IsNotExist(err) {
// 			return false
// 		}
// 	}
// 	return true
// }
// func uploadMedia(file multipart.File, filename string) {
// 	defer file.Close()
// 	tmpfile, _ := os.Create("../SOURCE/" + filename)
// 	defer tmpfile.Close()
// 	io.Copy(tmpfile, file)
// }

// func getMetadata(r *http.Request) ([]byte, error) {
// 	f, _, err := r.FormFile("metadata")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get metadata form file: %v", err)
// 	}

// 	metadata, errRead := ioutil.ReadAll(f)
// 	if errRead != nil {
// 		return nil, fmt.Errorf("failed to read metadata: %v", errRead)
// 	}

// 	return metadata, nil
// }

// func verifyRequest(r *http.Request) error {
// 	if _, ok := r.MultipartForm.File["media"]; !ok {
// 		return fmt.Errorf("media is absent")
// 	}

// 	if _, ok := r.MultipartForm.File["metadata"]; !ok {
// 		return fmt.Errorf("metadata is absent")
// 	}

// 	return nil
// }

func FilePathWalkDir(root string) ([]os.FileInfo, []string, error) {
	var files []os.FileInfo
	var paths []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			files = append(files, info)
			paths = append(paths, path)
		}
		return nil
	})
	return files, paths, err
}

func DeployFiles(c *gin.Context) {
	// if len(positionalArgs) == 0 {
	// 	log.Fatalf("This program requires at least 1 positional argument.")
	// }

	// Metadata content.

	// New multipart writer.
	var deployMeta models.Deployment
	if err := c.Bind(&deployMeta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Metadata part.
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json")
	metadataHeader.Set("Content-ID", "metadata")
	part, err := writer.CreatePart(metadataHeader)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Error writing metadata headers",
		})
		return
	}
	// files, readerr := ioutil.ReadDir("../SOURCE/go/")
	f, paths, readerr := FilePathWalkDir(deployMeta.SourceLink)
	for _, h := range f {
		fmt.Println(h.Name())
	}
	for _, p := range paths {
		fmt.Println(p)
	}
	metadata := deployMeta.DestinationLink
	part.Write([]byte(metadata))

	log.Println(f)
	if readerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "cannot read source",
		})
		return
	}
	// Media Files.
	// var meta models.MetaData
	// var paths []string
	for i, mediaFilename := range f {
		mediaData, errRead := ioutil.ReadFile(paths[i])
		// fmt.Println(mediaFilename)
		if errRead != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "Error reading media file",
			})
			return
		}
		// paths = append(paths, "../SOURCE/go/"+mediaFilename.Name())
		// println(paths)
		fileName := filepath.Base(mediaFilename.Name())
		log.Println(fileName)

		mediaHeader := textproto.MIMEHeader{}
		mediaHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\".", fileName))
		mediaHeader.Set("Content-ID", "media")
		mediaHeader.Set("Content-Filename", fileName)
		mediaHeader.Set("Content-Ticket", deployMeta.Ticket)
		mediaHeader.Set("Content-Filepath", paths[i])

		mediaPart, err := writer.CreatePart(mediaHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Error writing media headers",
				"err":    errRead,
			})
			return
		}

		if _, err := io.Copy(mediaPart, bytes.NewReader(mediaData)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Error writing media headers",
				"err":    errRead,
			})
			return
		}
	}

	// Close multipart writer.
	if err := writer.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "closing writer",
			"err":    err,
		})
		return
	}

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary())

	// Initialize HTTP Request and headers.
	uploadURL := "http://localhost:3002/api/deployMultiple"
	r, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error initializing a request",
			"err":    err,
		})
		return
	}
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Accept", "*/*")

	// HTTP Client.
	client := &http.Client{Timeout: 180 * time.Second}
	rsp, err := client.Do(r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error initializing a request",
			"err":    err,
		})
		return
	}

	fmt.Println(paths)
	// Check response status code.
	if rsp.StatusCode != http.StatusOK {

		c.JSON(rsp.StatusCode, gin.H{
			"status": "Request failed with response code",
			"err":    err,
		})
	} else {
		c.JSON(rsp.StatusCode, gin.H{
			"status": "Successfully Deployed",
		})
	}
}
