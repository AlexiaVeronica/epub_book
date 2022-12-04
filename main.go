package main

import (
	"epubset/pkg/config"
	"epubset/pkg/go-epub"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var Args *config.Config

func CreateFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err = os.Mkdir(filename, os.ModePerm); err != nil {
			fmt.Println("Mkdir error", err)
		}
	}
}
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
	epub *epub.Epub
}

func SetBookInfo(Author, Cover, Description string) *EpubConfig {
	Epub := &EpubConfig{epub: epub.NewEpub(Args.BookName)} // Create a new EPUB
	Epub.epub.SetLang("zh")
	Epub.epub.SetAuthor(Author)
	Epub.epub.SetDescription(Description)
	Epub.epub.SetCover(Cover, "")
	return Epub
}
func (ep *EpubConfig) DownloaderCover(CoverUrl string, Cover bool) {
	CreateFile("cover")
	resp, err := http.Get(CoverUrl)
	if err != nil {
		fmt.Println("Download error", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("Close error", err)
		}
	}(resp.Body)

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return
	}
	if content != nil {
		var coverName string
		if strings.Contains(path.Base(CoverUrl), ".jpg") {
			coverName = path.Join("cover", path.Base(CoverUrl))
		} else {
			coverName = path.Join("cover", path.Base(CoverUrl)+".jpg")
		}
		if err = os.WriteFile(coverName, content, 0666); err != nil {
			fmt.Println("WriteFile error", err)
		}
		image, ok := ep.epub.AddImage(coverName, "cover.jpg")
		if ok != nil {
			fmt.Println("AddImage error", ok)
		} else {
			fmt.Println("===>", image, "added")
			if Cover {
				ep.epub.SetCover(image, "")
			} else {
				ep.epub.AddSection("<img src=\""+image+"\"/>", "插画", "", "")
			}
		}
	} else {

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
func (ep *EpubConfig) Save() {
	// 判断文件夹是否存在
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err = os.Mkdir("output", os.ModePerm); err != nil {
			fmt.Println("Mkdir error", err)
			return
		}
	}
	err := ep.epub.Write(path.Join("output", Args.BookName+".epub"))
	if err != nil {
		fmt.Println("Save error", err)
	}
}

func (ep *EpubConfig) SplitChapter(file []byte) {
	var title string
	var content string
	title = "前言\n"
	for _, line := range strings.Split(string(file), "\n") {
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
	if file, err := os.ReadFile(Args.FileName); err != nil {
		fmt.Println("ReadFile error", err)
	} else {
		if Args.Cover != "" {
			Epub.DownloaderCover(Args.Cover, true)
			fmt.Println("===>", Args.Cover, "downloaded")
		}
		Epub.SplitChapter(file)
		Epub.Save()
	}

}
