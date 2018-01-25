package render

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
