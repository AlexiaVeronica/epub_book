package image

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func DownloaderCover(CoverUrl string, CoverName string) string {
	if CoverName == "" {
		CoverName = path.Base(CoverUrl)
	}
	CoverName = path.Join("cover", CoverName)
	if !strings.Contains(CoverName, ".jpg") {
		CoverName = CoverName + ".jpg"
	}
	// 判断图片是否存在
	_, exist := os.Stat(CoverName)
	if exist == nil {
		return CoverName
	}
	resp, err := http.Get(CoverUrl)
	if err != nil {
		fmt.Println("Download error", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			fmt.Println("Close error", err)
		}
	}(resp.Body)

	if imageFile, err := io.ReadAll(resp.Body); err != nil {
		fmt.Printf("Read http response failed! %v", err)
	} else {
		if ok := os.WriteFile(CoverName, imageFile, 0666); ok != nil {
			fmt.Println("WriteFile error", ok)
		}
		return CoverName
	}
	return ""
}
