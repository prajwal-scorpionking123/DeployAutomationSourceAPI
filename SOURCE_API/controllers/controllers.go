package controllers

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/team_six/SOURCE_API/helpers"
	"github.com/team_six/SOURCE_API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database = helpers.ConnectToMongoDB()
var collection *mongo.Collection = DB.Collection("Source")

//checking whether file in that entered path is exist or not
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func PostLink(c *gin.Context) {

	c.Header("content-type", "application/json")

	var source models.Source

	if err := c.Bind(&source); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	source.Timestamp = time.Now()
	if !FileExists(source.SourceLink) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "there is no file entered path"})
		return
	}
	res, err := collection.InsertOne(context.TODO(), source)
	if err != nil {
		helpers.GetError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"source": res})
}
func GetSources(c *gin.Context) {
	c.Header("content-type", "application/json")
	var sources []models.Source
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helpers.GetError(err, c)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var source models.Source
		err := cur.Decode(&source)
		if err != nil {
			log.Fatal(err)
		}

		sources = append(sources, source)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"sources": sources})
}

func DeployFiles(c *gin.Context) {
	var deployMeta models.Deployment
	if err := c.Bind(&deployMeta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//creating client object
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	// c.FileAttachment("./assets/go/m.go", "m.go")

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	//opering the source file to be deplyed
	file, err := os.Open(deployMeta.SourceLink)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "source file not found",
		})
		return

	}
	fileName := filepath.Base(file.Name())
	println(fileName)

	// creating the form field
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "failed",
		})
		return

	}
	_, err = io.Copy(fw, file)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "failed",
		})
		return

	}
	writer.Close()

	//calling the production api to deploy the soruce file
	req, err := http.NewRequest("POST", "http://localhost:3002/api/deployFiles", bytes.NewReader(body.Bytes()))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "failed",
		})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rsp, _ := client.Do(req)

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
		c.JSON(rsp.StatusCode, gin.H{
			"Request failed with response code": rsp.Status,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "OK",
	})
}
