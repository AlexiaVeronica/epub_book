package main

import (
	"epubset/go-epub"
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

func SetBookInfo(Author, Cover, Description, Title string) {
	ep.SetLang("zh")
	ep.SetAuthor(Author)
	ep.SetTitle(Title)
	ep.SetDescription(Description)
	ep.SetCover(Cover, "")
}
func Save() {
	err := ep.Write(book_name + ".epub")
	if err != nil {
		fmt.Println("Save error", err)
	}
}

const DefaultMatchTips = "^第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^\\d{1,4}$|^引子$|^楔子$|^章节目录|^章节|^序章"

func main() {
	// Create a new EPUB
	file_name := "DIO异闻录.txt"
	book_name = strings.Replace(file_name, ".txt", "", -1)
	ep = epub.NewEpub("My title")
	SetBookInfo("Author", "cover.jpg", "Description", "Title")
	file, _ := os.ReadFile(file_name)
	//fmt.Println(string(file))

	reg, err := regexp.Compile(DefaultMatchTips)
	if err != nil {
		fmt.Println("regexp error", err)
	} else {
		fmt.Println("regexp ok")

	}
	var title string
	var content string
	title = "前言"
	for _, line := range strings.Split(string(file), "\n") {
		if reg.MatchString(line) {
			AddChapter(title, "<h1>"+strings.ReplaceAll(title, "\r", "")+"</h1>"+content)
			title = line
			content = ""
		} else {
			content += fmt.Sprintf("\n<p>%s</p>", strings.ReplaceAll(line, "\r", ""))
		}
	} //end for
	Save()
}
