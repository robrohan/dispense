package models

// Config is the config object for the application
type Config struct {
	Base struct {
		Input  string `conf:"default:./test_data"`
		Output string `conf:"default:./build"`
	}
	Template struct {
		Extension string `conf:"default:tpl"`
		Directory string `conf:"default:./public"`
		Listing   string `conf:"default:index"`
		Post      string `conf:"default:post"`
	}
}
