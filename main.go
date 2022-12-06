package main

import (
	"epubset/pkg/config"
	"epubset/pkg/epub"
	"epubset/pkg/file"
	"epubset/pkg/image"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

var Args *config.Config

func init() {
	Args = config.InitParams()
	if Args.FileName == "" {
		fmt.Println("Please input file name, use -h to get help")
		os.Exit(0)
	}
	if Args.BookName == "" {
		Args.BookName = strings.ReplaceAll(Args.FileName, ".txt", "")
	}

}

type EpubConfig struct {
	// Epub is the main struct for the epub package.
	epub     *epub.Epub
	saveDir  string
	savePath string
}

func SetBookInfo(Author, Cover, Description string) *EpubConfig {
	Epub := &EpubConfig{epub: epub.NewEpub(Args.BookName), saveDir: "output"} // Create a new EPUB
	if Args.SaveDir != "" {
		Epub.saveDir = Args.SaveDir
	}

	file.CreateFile(Epub.saveDir)
	file.CreateFile("cover")
	Epub.savePath = path.Join(Epub.saveDir, Args.BookName+".epub") // Set the output file path
	Epub.epub.SetLang("zh")
	Epub.epub.SetAuthor(Author)
	Epub.epub.SetDescription(Description)
	Epub.epub.SetCover(Cover, "")
	return Epub
}

func (ep *EpubConfig) DownloaderCover(CoverUrl string, Cover bool) {
	coverName := image.DownloaderCover(CoverUrl)
	FilePath, ok := ep.epub.AddImage(coverName, "cover.jpg")
	if ok != nil {
		fmt.Println("AddImage error", ok)
	} else {
		fmt.Println("===>", FilePath, "added")
		if Cover {
			ep.epub.SetCover(FilePath, "")
		} else {
			ep.AddChapter("封面", fmt.Sprintf(`<img src="%s" />`, FilePath))
		}
	}
}

func (ep *EpubConfig) AddChapter(title string, content string) {
	_, err := ep.epub.AddSection(content, title, "", "")
	if err != nil {
		fmt.Println("AddSection error", err)
		return
	}
	//println(section) // section0002.xhtml
}
func (ep *EpubConfig) Save(max_retry int) {
	if err := ep.epub.Write(ep.savePath); err != nil {
		fmt.Println("Save error:", err)
		if max_retry > 0 {
			ep.Save(max_retry - 1)
		}
	}
}

func (ep *EpubConfig) SplitChapter(fileByte []byte) {
	var title string
	var content string
	title = "前言\n"
	for _, line := range strings.Split(string(fileByte), "\n") {
		line = strings.ReplaceAll(line, "\r", "")
		if regexp.MustCompile(Args.Rule).MatchString(line) {
			ep.AddChapter(title, "<h1>"+title+"</h1>"+content) // 添加章节
			title = line                                       // new title
			content = ""                                       // clear content
		} else {
			content += fmt.Sprintf("\n<p>%s</p>", line)
		}
	} //end for
	fmt.Println(Args.BookName, "done") // last chapter
}

func main() {
	Epub := SetBookInfo(Args.Author, Args.Cover, Args.Description)
	if fileByte, err := os.ReadFile(Args.FileName); err != nil {
		fmt.Println("ReadFile error", err)
	} else {
		if Args.Cover != "" {
			Epub.DownloaderCover(Args.Cover, true)
			fmt.Println("===>", Args.Cover, "downloaded")
		}
		Epub.SplitChapter(fileByte)
		Epub.Save(2)
	}

}
