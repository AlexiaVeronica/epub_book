package config

import "flag"

const MatchTips = "^第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^\\d{1,4}$|^引子$|^楔子$|^章节目录|^章节|^序章"

var file_name = flag.String("n", "", "file_name")
var author = flag.String("a", "", "author name")
var cover = flag.String("c", "", "cover image")
var description = flag.String("d", "", "description")
var rule = flag.String("r", MatchTips, "rule")

type Config struct {
	FileName    string
	Author      string
	Cover       string
	Description string
	Rule        string
	BookName    string
}

func InitParams() *Config {
	flag.Parse()
	return &Config{
		FileName:    *file_name,
		Author:      *author,
		Cover:       *cover,
		Description: *description,
		Rule:        *rule,
	}
}
