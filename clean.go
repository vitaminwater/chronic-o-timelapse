package main

import (
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
	//name := MustGetenv("NAME")

	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)
	path := "" //fmt.Sprintf("/%s", name)
	res, err := dbx.ListFolder(files.NewListFolderArg(path))
	fu(err)

	for _, m := range res.Entries {
		switch f := m.(type) {
		case *files.FileMetadata:
			_, err = dbx.DeleteV2(files.NewDeleteArg(f.PathLower))
			fu(err)
			logrus.Infof("Deleted file: %s", f.PathLower)
		case *files.FolderMetadata:
			logrus.Infof("Folder: %s", f.Name)
			_, err = dbx.DeleteV2(files.NewDeleteArg(f.PathLower))
			fu(err)
			logrus.Infof("Deleted folder: %s", f.PathLower)
		case *files.DeletedMetadata:
			logrus.Infof("Deleted: %s", f.Name)
		}
	}
}
