package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var dbx files.Client

func init() {
	token := MustGetenv("DBX_TOKEN")
	config := dropbox.Config{
		Token: token,
	}
	dbx = files.New(config)
}

func fu(e error) {
	if e != nil {
		logrus.Fatal(e)
	}
}

func MustGetenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logrus.Fatalf("missing env %s", name)
	}
	return v
}

func GetFileReader(name string) (io.ReadCloser, error) {
	dbxFile := fmt.Sprintf("/%s/latest.jpg", name)

	da := files.NewDownloadArg(dbxFile)

	_, c, err := dbx.Download(da)

	logrus.Infof("Downloaded %s", dbxFile)
	return c, err
}

func serve(c *gin.Context) {
	name := c.Param("name")

	content, err := GetFileReader(name)
	if err != nil {
		logrus.Warning(err)
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "image/jpeg")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, content)
}

func main() {
	r := gin.Default()
	r.GET("/:name", serve)

	certFile := "certs/chronic-o-matic.com.crt"
	keyFile := "certs/chronic-o-matic.com.key"
	if _, err := os.Stat(certFile); err == nil {
		logrus.Info("exists1")
		if _, err := os.Stat(keyFile); err == nil {
			logrus.Info("exists2")
			go r.RunTLS(":443", certFile, keyFile)
		}
	}
	r.Run(":80")
}
