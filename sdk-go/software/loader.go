package software

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/richardlt/matrix/internal/errors"
)

var is []*Image
var ts []*Theme
var fs []*Font

type file struct {
	Name string
	Data []byte
}

func loadImages() error {
	files, err := loadFiles("images")
	if err != nil {
		return errors.Errorf("loading images files: %w", err)
	}

	for _, file := range files {
		i := &Image{}
		if err := json.Unmarshal(file.Data, i); err != nil {
			return errors.Errorf("unmarshaling image file %s: %w", file.Name, err)
		}
		is = append(is, i)
	}

	return nil
}

func loadThemes() error {
	files, err := loadFiles("themes")
	if err != nil {
		return errors.Errorf("loading themes files: %w", err)
	}

	for _, file := range files {
		t := &Theme{}
		if err := json.Unmarshal(file.Data, t); err != nil {
			return errors.Errorf("unmarshaling theme file %s: %w", file.Name, err)
		}
		ts = append(ts, t)
	}

	return nil
}

func loadFonts() error {
	files, err := loadFiles("fonts")
	if err != nil {
		return errors.Errorf("loading fonts files: %w", err)
	}

	for _, file := range files {
		f := &Font{}
		if err := json.Unmarshal(file.Data, f); err != nil {
			return errors.Errorf("unmarshaling font file %s: %w", file.Name, err)
		}
		fs = append(fs, f)
	}

	return nil
}

func loadFiles(dir string) ([]file, error) {
	files, err := os.ReadDir("./" + dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Errorf("reading directory %s: %w", dir, err)
	}

	res := []file{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		data, err := os.ReadFile(fmt.Sprintf("./%s/%s", dir, f.Name()))
		if err != nil {
			return nil, errors.Errorf("reading %s/%s: %w", dir, f.Name(), err)
		}

		res = append(res, file{f.Name(), data})
	}

	return res, nil
}
