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
	}

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)
	path := ""
	if name != "" {
		path = fmt.Sprintf("/%s", name)
	}
	args := files.NewListFolderArg(path)
	res, err := dbx.ListFolder(args)
	fu(err)

	for {
		for _, m := range res.Entries {
			switch f := m.(type) {
			case *files.FileMetadata:
				logrus.Infof("f %s", f.PathLower)
			case *files.FolderMetadata:
				logrus.Infof("d %s/", f.Name)
			case *files.DeletedMetadata:
				logrus.Infof("- %s", f.Name)
			}
		}
		if res.HasMore {
			_, err = dbx.ListFolderContinue(files.NewListFolderContinueArg(res.Cursor))
			fu(err)
		} else {
			break
		}
	}
}
