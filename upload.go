package main

import (
	"fmt"
	"os"

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
	name := ""
	if len(os.Args) >= 2 {
		name = os.Args[1]
	} else {
		logrus.Fatal("missing arg name")
	}

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)

	local := "image_00001.jpg"
	remote := "lol.jpg"

	f, err := os.Open(local)
	fu(err)

	ci := files.NewCommitInfo(fmt.Sprintf("/%s/%s", name, remote))
	res, err := dbx.Upload(ci, f)
	fu(err)

	logrus.Info(res)
}
