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
	t := time.Now().Unix()
	name := fmt.Sprintf("%d.jpg", t)
	cmd := exec.Command("/usr/bin/raspistill", "-o", name)
	err := cmd.Run()
	return name, err
}

func main() {
	token := MustGetenv("DBX_TOKEN")
	name := MustGetenv("NAME")

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)

	local, err := takePic()
	fu(err)

	remote := fmt.Sprintf("%d.jpg", int32(time.Now().Unix()))

	f, err := os.Open(local)
	fu(err)

	p := fmt.Sprintf("/%s/%s", name, remote)
	ci := files.NewCommitInfo(p)
	_, err = dbx.Upload(ci, f)
	fu(err)

	logrus.Infof("Uploaded %s", p)
}
