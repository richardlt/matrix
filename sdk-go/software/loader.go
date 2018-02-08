package software

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

var is []Image
var ts []Theme
var fs []Font

type file struct {
	Name string
	Data []byte
}

func loadImages() error {
	files, err := loadFiles("images")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		var i Image
		if err := json.Unmarshal(file.Data, &i); err != nil {
			return fmt.Errorf("Can't unmarshal %s image file", file.Name)
		}
		is = append(is, i)
	}

	return nil
}

func loadThemes() error {
	files, err := loadFiles("themes")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		var t Theme
		if err := json.Unmarshal(file.Data, &t); err != nil {
			return fmt.Errorf("Can't unmarshal %s theme file", file.Name)
		}
		ts = append(ts, t)
	}

	return nil
}

func loadFonts() error {
	files, err := loadFiles("fonts")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		var f Font
		if err := json.Unmarshal(file.Data, &f); err != nil {
			return fmt.Errorf("Can't unmarshal %s font file", file.Name)
		}
		fs = append(fs, f)
	}

	return nil
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
