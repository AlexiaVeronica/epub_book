package config

import "flag"

const MatchTips = "^第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^\\d{1,4}$|^引子$|^楔子$|^章节目录|^章节|^序章"

var file = flag.String("file", "", "Input file name (required)")
var author = flag.String("author", "", "author name")
var cover = flag.String("cover", "", "cover image")
var description = flag.String("intro", "", "description information")

var transformTW = flag.Bool("tw", false, "Transform to traditional Chinese")
var transformZh = flag.Bool("zh", false, "Transform to simplified Chinese")

var rule = flag.String("r", MatchTips, "rule")
var saveDir = flag.String("o", "", "output dir")
var coverName = flag.String("cover_name", "", "set cover name")
var showProgress = flag.Bool("p", false, "show progress bar")

type Config struct {
	FileName     string
	Author       string
	Cover        string
	Description  string
	Rule         string
	BookName     string
	SaveDir      string
	CoverName    string
	TransformTW  bool
	TransformZh  bool
	ShowProgress bool
}

func InitParams() *Config {
	flag.Parse()
	return &Config{
		FileName:     *file,
		Author:       *author,
		Cover:        *cover,
		Description:  *description,
		Rule:         *rule,
		TransformTW:  *transformTW,
		TransformZh:  *transformZh,
		SaveDir:      *saveDir,
		ShowProgress: *showProgress,
		CoverName:    *coverName,
	}
}
