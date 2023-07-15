package models

type Config struct {
	Base struct {
		// Input is where the markdown data lives. The files we want to put on the internet
		Input string `conf:"default:./test_data"`
		// Output is the directory where the rendered HTML will be written - the dir to upload to the sever
		Output string `conf:"default:./public"`
	}
	Template struct {
		Extension string `conf:"default:tpl"`
		Directory string `conf:"default:./hive"`
		Listing   string `conf:"default:index"`
		Post      string `conf:"default:post"`
	}
}
