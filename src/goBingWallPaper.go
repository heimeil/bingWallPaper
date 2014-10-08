package main

import (
	"syscall"
	"unsafe"
	"strings"
	"path/filepath"
	"os"
	"os/exec"
	"regexp"
	"net/http"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	fpath := GetBingWallPaper()
	if fpath != "" {
		SetWallPaper(fpath)
	}else {
		log.Fatalln("not find img")
	}
}

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}

func GetBingWallPaper() string {
	reg, err := regexp.Compile(`http://s.cn.bing.net/az/hprichbg/rb/[\w\W]*?\.\w{3,4}`)
	CheckErr(err)
	body := Get("http://cn.bing.com")
	str := reg.FindString(string(body))
	if str != "" {
		body := Get(str)
		err := os.MkdirAll(GetCurrPath()+"bing/", os.ModePerm)
		CheckErr(err)
		fpath := GetCurrPath() + "bing/" + time.Now().Format("20060102") + str[strings.LastIndex(str, "."):]
		err = ioutil.WriteFile(fpath, body, os.ModePerm)
		CheckErr(err)
		return fpath
	}
	return ""
}

func SetWallPaper(fpath string) {
	h := syscall.MustLoadDLL("user32.dll")
	c := h.MustFindProc("SystemParametersInfoW")
	defer syscall.FreeLibrary(h.Handle)
	uiAction := 0x0014
	uiParam := 0
	pvParam := syscall.StringToUTF16Ptr(fpath)
	fWinIni := 1
	r2, _, err := c.Call(uintptr(uiAction),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(pvParam)),
		uintptr(fWinIni))
	if r2 != 0 {
		log.Fatalln(r2, err, fpath)
	}
}

func Get(url string) []byte {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	return body
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
