package render

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Init prepares rendering stuff.
func Init() error {
	if err := loadImages(); err != nil {
		return err
	}
	if err := loadThemes(); err != nil {
		return err
	}
	return loadFonts()
}

type file struct {
	Name string
	Data []byte
}

func loadFiles(dir string) ([]file, error) {
	files, err := ioutil.ReadDir("./" + dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	res := []file{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		data, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s", dir, f.Name()))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		res = append(res, file{f.Name(), data})
	}

	return res, nil
}
