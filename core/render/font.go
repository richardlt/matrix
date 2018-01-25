package render

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Caracter is a pixel template for a given caracter.
type Caracter struct {
	Width uint64   `json:"width"`
	Mask  []uint64 `json:"mask"`
}

// NewFont returns a new font.
func NewFont() Font { return Font{} }

// Font contains pixel templates for caracters.
type Font struct {
	Name      string              `json:"name"`
	Height    uint64              `json:"height"`
	Caracters map[string]Caracter `json:"caracters"`
}

// GetCaracterByValue returns a caracter from given value.
func (f *Font) GetCaracterByValue(value rune) Caracter {
	for k, c := range f.Caracters {
		if []rune(k)[0] == value {
			return c
		}
	}

	return Caracter{}
}

// GetFontByName returns an loaded font in memory.
func GetFontByName(name string) Font {
	for _, f := range fs {
		if f.Name == name {
			return f
		}
	}
	return Font{}
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
