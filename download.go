package main

import (
	"io"
	"os"
	"path"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/sirupsen/logrus"
)

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

func main() {
	token := MustGetenv("DBX_TOKEN")
	dbxFile := ""
	if len(os.Args) >= 2 {
		dbxFile = os.Args[1]
	} else {
		logrus.Fatal("missing arg dbxFile")
	}

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)

	da := files.NewDownloadArg(dbxFile)

	_, c, err := dbx.Download(da)
	fu(err)

	outFileName := path.Base(dbxFile)
	outFile, err := os.Create(outFileName)
	fu(err)
	defer outFile.Close()

	_, err = io.Copy(outFile, c)
	logrus.Infof("Downloaded %s", outFileName)
}
