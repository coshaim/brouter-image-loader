package downloader

// Configuration of downloader
type Configuration struct {
	Segments4URL string `env:"SEGMENTS4_URL"`
	FilePath     string `env:"FILE_PATH"`
}
