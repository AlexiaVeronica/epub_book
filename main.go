package main

import (
	"epubset/go-epub"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func AddChapter(title string, content string) {
	_, err := ep.AddSection(content, title, "", "")
	if err != nil {
		return
	}
	//println(section) // section0002.xhtml
}

var ep *epub.Epub
var book_name string

func SetBookInfo(Author, Cover, Description string) {
	ep.SetLang("zh")
	ep.SetAuthor(Author)
	ep.SetDescription(Description)
	ep.SetCover(Cover, "")
}
func Save() {
	err := ep.Write(book_name + ".epub")
	if err != nil {
		fmt.Println("Save error", err)
	}
}

func SplitChapter(file []byte) {
	var title string
	var content string
	title = "前言\n"
	for _, line := range strings.Split(string(file), "\n") {
		if regexp.MustCompile(MatchTips).MatchString(line) {
			AddChapter(title, "<h1>"+strings.ReplaceAll(title, "\r", "")+"</h1>"+content)
			title = line
			content = ""
		} else {
			content += fmt.Sprintf("\n<p>%s</p>", strings.ReplaceAll(line, "\r", ""))
		}
	} //end for
	fmt.Println("Done")
}

const MatchTips = "^第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^\\d{1,4}$|^引子$|^楔子$|^章节目录|^章节|^序章"

// var b = flag.Bool("b", false, "bool类型参数")
var file_name = flag.String("n", "", "string类型参数")

func main() {
	flag.Parse()
	if *file_name == "" {
		fmt.Println("Usage: -n book_name")
		return
	}
	book_name = strings.Replace(*file_name, ".txt", "", -1)
	ep = epub.NewEpub(book_name) // Create a new EPUB
	SetBookInfo("Author", "cover.jpg", "Description")
	file, err := os.ReadFile(*file_name)
	if err != nil {
		fmt.Println("ReadFile error", err)
		return
	} else {
		SplitChapter(file)
		Save()
	}
}
