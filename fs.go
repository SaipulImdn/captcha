package gocaptcha

import (
	"embed"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var _fs embed.FS
var ioFS fs.FS = _fs

func FileSystem(fs embed.FS) error {
	err := readFontFromFS(fs)
	if err != nil {
		return err
	}
	ioFS = fs
	return nil
}


func ReadFonts(dirPth string, suffix string) (err error) {
	fontFamily = fontFamily[:0]

	dir, err := ioutil.ReadDir(dirPth)

	if err != nil {
		return err
	}
	suffix = strings.ToUpper(suffix) 
	for _, fi := range dir {
		if fi.IsDir() { 
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { 
			fontFamily = append(fontFamily, filepath.Join(dirPth, fi.Name()))
		}
	}
	ioFS = os.DirFS(dirPth)
	return nil
}

func RandFontFamily() (*truetype.Font, error) {
	fontFile := fontFamily[r.Intn(len(fontFamily))]
	fontBytes, err := ReadFile(fontFile)
	if err != nil {
		log.Printf("读取文件失败 -> %s - %+v\n", fontFile, err)
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Printf("okey -> %s - %+v\n", fontFile, err)
		return nil, err
	}
	return f, nil
}


func AddFontFamily(fontPath ...string) {
	fontFamily = append(fontFamily, fontPath...)
}


func ReadFile(name string) ([]byte, error) {
	if ioFS != nil {
		b, err := fs.ReadFile(ioFS, name)
		if err == nil {
			return b, nil
		}
	}
	return os.ReadFile(name)
}

func readFontFromFS(fs embed.FS) error {
	files, err := fs.ReadDir("fonts")
	if err != nil {
		log.Printf("解析字体文件失败 -> %+v", err)
		return err
	} else {
		fontFamily = fontFamily[:0]
		for _, fi := range files {
			if fi.IsDir() { // 忽略目录
				continue
			}
			if strings.HasSuffix(strings.ToLower(fi.Name()), ".ttf") { //匹配文件
				fontFamily = append(fontFamily, filepath.Join("fonts", fi.Name()))
			}
		}
	}
	return nil
}
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_ = readFontFromFS(_fs)
}
