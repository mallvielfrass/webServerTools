package wst

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-chi/chi"
	"github.com/mallvielfrass/fmc"
)

func CheckAccessArea(path string) (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	compare := strings.Contains(abs, dir)
	//fmt.Println("PWD:", dir)
	//fmt.Println("Absolute:", abs)
	ex := fileExists(abs)
	//fmt.Printf("Contains: %t | FileExists: %t\n", compare, ex)
	if !compare || !ex {
		return "", false
	}
	return abs, true
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func GetType(path string) string {
	s := strings.Split(path, "/")
	if len(s) == 0 {
		return "undefined"
	}
	file := s[len(s)-1]
	f := strings.Split(file, ".")
	if len(f) == 0 {
		return "undefined"
	}
	extension := f[len(f)-1]
	return extension
}
func staticRouter(w http.ResponseWriter, r *http.Request) {
	urlFile := r.URL.Path
	ext := GetType(urlFile)
	info, area := CheckAccessArea("." + urlFile)
	if !area {
		fmc.Printfln("#bbtStaticRouter> #rbtError: #ybtURL not in access area:#bbt[#gbt%s#bbt]", urlFile)
		http.NotFound(w, r)
		return
	}
	fmc.Printfln("#bbtStaticRouter> #ybtURL:#bbt[#gbt%s#bbt] #ybtType:#bbt[#gbt%s#bbt] #ybtLocal File:#bbt[#gbt%s#bbt]", urlFile, ext, info)
	switch ext {
	case "css":
		w.Header().Set("Content-Type", "text/css")
		fmc.Printfln("#bbtStaticRouter> #ybtSet header: #bbt[#gbt%s#bbt]", ext)
	case "png":
		w.Header().Set("Content-Type", "image/png")
		fmc.Printfln("#bbtStaticRouter> #ybtSet header: #bbt[#gbt%s#bbt]", ext)
	case "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		fmc.Printfln("#bbtStaticRouter> #ybtSet header: #bbt[#gbt%s#bbt]", ext)
	case "js":
		w.Header().Set("Content-Type", "application/javascript")
		fmc.Printfln("#bbtStaticRouter> #ybtSet header: #bbt[#gbt%s#bbt]", ext)
	case "ttf":
		w.Header().Set("Content-Type", "application/x-font-ttf")
		fmc.Printfln("#bbtStaticRouter> #ybtSet header: #bbt[#gbt%s#bbt]", ext)
	default:
		fmc.Printfln("#bbtStaticRouter> Undefined type [%s] of file: [%s]", ext, urlFile)
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, info)
}
func (folder staticFolder) StaticRouter(w http.ResponseWriter, r *http.Request) {
	file := chi.URLParam(r, "file")
	typeFile := chi.URLParam(r, "type")
	switch typeFile {
	case "css":
		log.Printf("Type [%s] of file: [%s]\n", typeFile, file)
		w.Header().Set("Content-Type", "text/css")
	case "js":
		log.Printf("Type [%s] of file: [%s]\n", typeFile, file)
		w.Header().Set("Content-Type", "application/javascript")
	case "ttf":
		log.Printf("Type [%s] of file: [%s]\n", typeFile, file)
		w.Header().Set("Content-Type", "application/x-font-ttf")
	default:
		log.Printf("Undefined type [%s] of file: [%s]\n", typeFile, file)
	}
	path := "./web/static/" + typeFile + "/" + file
	//	log.Println(path)
	fmt.Fprint(w, path)
	//http.ServeFile(w, r, path)
}
func OpenFile(file string) []byte {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(file)
	fmc.ErrorHandleFatal(err, "Open file: "+file)
	io.Copy(buf, f)
	f.Close()
	return buf.Bytes()
}
func (folder staticFolder) fileRouter(w http.ResponseWriter, r *http.Request) {
	//fmc.Printfln("folde: %s", folder.Path)
	path := r.URL.Path
	//	fmt.Println(path)
	//	fmt.Fprint(w, r.URL)
	norm := strings.SplitAfter(path, "/static/")[1]
	//	fmc.Println(norm)
	file := folder.Path + norm
	p, ok := CheckAccessArea(file)
	if !ok {
		fmc.Printfln("[#ybtFileRouter#RRR]:\n\t#ybtWarning#RRR: Access denied to file #RRR[#gbt%s#RRR]", file)
		fmt.Fprint(w, http.StatusNotFound)
	} else {
		//http.
		//fmt.Println(p)
		mim, err := mimetype.DetectFile(p)
		fmc.ErrorHandle(err)
		//fmt.Println(mim.String())
		//w.Header().Set("Content-Type", mim.String())

		Openfile, err := os.Open(p)
		if err != nil {
			//File not found, send 404
			http.Error(w, "File not found.", 404)
			return
		}
		defer Openfile.Close() //Close after function return
		//File is found, create and send the correct headers

		//Get the Content-Type of the file
		//Create a buffer to store the header of the file in
		FileHeader := make([]byte, 512)
		//Copy the headers into the FileHeader buffer
		Openfile.Read(FileHeader)
		//Get content type of file
		//FileContentType := http.DetectContentType(FileHeader)

		//Get the file size
		FileStat, _ := Openfile.Stat()                     //Get info from file
		FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string
		nameSplit := strings.Split(p, "/")
		name := nameSplit[len(nameSplit)-1]
		//Send the headers
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		w.Header().Set("Content-Type", mim.String())
		w.Header().Set("Content-Length", FileSize)
		fmc.Printfln("[#ybtFileRouter#RRR]:\n\t#bbtAccess#RRR: #RRR Access is allowed to file #RRR[#gbt%s#RRR] \n\t#bbtName#RRR: #RRR[#gbt%s#RRR]\n\t#bbtMimtype#RRR: #RRR[#gbt%s#RRR]\n\t#bbtSize#RRR: #RRR[#gbt%s bytes#RRR]", file, name, mim.String(), FileSize)

		//Send the file
		//We read 512 bytes from the file already, so we reset the offset back to 0
		Openfile.Seek(0, 0)
		io.Copy(w, Openfile)
	}

}
func checkFolder(path string) string {
	if 2 < len(path) {
		if path[len(path)-1] != '/' {
			path = path + "/"
		}
		if path[0:2] != "./" {
			if path[0:1] != "/" {
				path = "./" + path
			}

		}
		return path
	}
	return "./"
}
func FileServer(r chi.Router, folder string) {

	r.HandleFunc("/static/*", staticFolder{checkFolder(folder)}.fileRouter)
}
