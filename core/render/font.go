package render

import (
	"encoding/json"
	"fmt"

	"github.com/richardlt/matrix/sdk-go/software"

	"github.com/pkg/errors"
)

// GetFontCaracterByValue returns a font's caracter from given value.
func GetFontCaracterByValue(f software.Font, value rune) software.Font_Caracter {
	for k, c := range f.Caracters {
		if []rune(k)[0] == value {
			return *c
		}
	}

	return software.Font_Caracter{}
}

// GetFontByName returns an loaded font in memory.
func GetFontByName(name string) software.Font {
	for _, f := range fs {
		if f.Name == name {
			return f
		}
	}
	return software.Font{}
}

var fs []software.Font

func loadFonts() error {
	files, err := loadFiles("fonts")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		var f software.Font
		if err := json.Unmarshal(file.Data, &f); err != nil {
			return fmt.Errorf("Can't unmarshal %s font file", file.Name)
		}
		fs = append(fs, f)
	}

	return nil
}
