package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
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

func takePic() (string, error) {
	name := "latest.jpg"
	cmd := exec.Command("/usr/bin/raspistill", "-vf", "-hf", "-q", "50", "-o", name)
	err := cmd.Run()
	return name, err
}

func uploadPic(name, local, remote string) {
	f, err := os.Open(local)
	fu(err)

	p := fmt.Sprintf("/%s/%s", name, remote)
	ci := files.NewCommitInfo(p)
	ci.Mode.Tag = "overwrite"
	_, err = dbx.Upload(ci, f)
	fu(err)

	logrus.Infof("Uploaded %s", p)
}

func main() {
	name := MustGetenv("NAME")

	remote := fmt.Sprintf("%d.jpg", int32(time.Now().Unix()))
	local, err := takePic()
	fu(err)

	logrus.Info("Uploading files")
	uploadPic(name, local, remote)
	uploadPic(name, local, "latest.jpg")
}
