package main

import (
	"epubset/pkg/config"
	"epubset/pkg/go-epub"
	"fmt"
	"os"
	"regexp"
	"strings"
)

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
	err := ep.Write(Args.BookName + ".epub")
	if err != nil {
		fmt.Println("Save error", err)
	}
}

func SplitChapter(file []byte) {
	var title string
	var content string
	title = "前言\n"
	for _, line := range strings.Split(string(file), "\n") {
		if regexp.MustCompile(Args.Rule).MatchString(line) {
			AddChapter(title, "<h1>"+strings.ReplaceAll(title, "\r", "")+"</h1>"+content)
			title = line
			content = ""
		} else {
			content += fmt.Sprintf("\n<p>%s</p>", strings.ReplaceAll(line, "\r", ""))
		}
	} //end for
	fmt.Println("Done")
}

var Args *config.Config

func init() {
	Args = config.InitParams()
	if Args.FileName == "" {
		fmt.Println("Please input file name")
		os.Exit(0)
	}
	Args.BookName = strings.ReplaceAll(Args.FileName, ".txt", "")

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
