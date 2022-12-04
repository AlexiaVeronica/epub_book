package main

import (
	"epubset/pkg/config"
	"epubset/pkg/go-epub"
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

func AddChapter(title string, content string) {
	_, err := ep.AddSection(content, title, "", "")
	if err != nil {
		fmt.Println("AddSection error", err)
		return
	}
	//println(section) // section0002.xhtml
}

var ep *epub.Epub

func SetBookInfo(Author, Cover, Description string) {
	ep.SetLang("zh")
	ep.SetAuthor(Author)
	ep.SetDescription(Description)
	ep.SetCover(Cover, "")
}
func Save() {
	// 判断文件夹是否存在
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		os.Mkdir("output", os.ModePerm)
	}
	err := ep.Write(path.Join("output", Args.BookName+".epub"))
	if err != nil {
		fmt.Println("Save error", err)
	}
}

func SplitChapter(file []byte) {
	var title string
	var content string
	title = "前言\n"
	for _, line := range strings.Split(string(file), "\n") {
		line = strings.ReplaceAll(line, "\r", "")
		if regexp.MustCompile(Args.Rule).MatchString(line) {
			AddChapter(title, "<h1>"+title+"</h1>"+content)
			title = line // new title
			content = "" // clear content
		} else {
			content += fmt.Sprintf("\n<p>%s</p>", line)
		}
	} //end for
	fmt.Println(Args.BookName, "done")
}

func main() {
	ep = epub.NewEpub(Args.BookName) // Create a new EPUB
	SetBookInfo(Args.Author, Args.Cover, Args.Description)
	if file, err := os.ReadFile(Args.FileName); err != nil {
		fmt.Println("ReadFile error", err)
	} else {
		SplitChapter(file)
		Save()
	}
}
