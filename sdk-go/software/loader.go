package software

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

var is []Image

func loadImages() error {
	files, err := ioutil.ReadDir("./images")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		file, err := ioutil.ReadFile("./images/" + file.Name())
		if err != nil {
			return errors.WithStack(err)
		}

		var i Image
		if err := json.Unmarshal(file, &i); err != nil {
			return errors.WithStack(err)
		}

		is = append(is, i)
	}

	return nil
}

var ts []Theme

func loadThemes() error {
	files, err := ioutil.ReadDir("./themes")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		file, err := ioutil.ReadFile("./themes/" + file.Name())
		if err != nil {
			return errors.WithStack(err)
		}

		var t Theme
		if err := json.Unmarshal(file, &t); err != nil {
			return errors.WithStack(err)
		}

		ts = append(ts, t)
	}

	return nil
}

var fs []Font

func loadFonts() error {
	files, err := ioutil.ReadDir("./fonts")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		file, err := ioutil.ReadFile("./fonts/" + file.Name())
		if err != nil {
			return errors.WithStack(err)
		}

		var f Font
		if err := json.Unmarshal(file, &f); err != nil {
			return errors.WithStack(err)
		}

		fs = append(fs, f)
	}

	return nil
}
