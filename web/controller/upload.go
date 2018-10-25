/**
  @Prject: goProjects
  @Dev_Software: GoLand
  @File : upload
  @Time : 2018/10/24 11:58 
  @Author : hanxiaodong
*/

package controller

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"crypto/rand"
	"path/filepath"
	"os"
	"mime"
	"log"
)


func (app *Application) UploadFile(w http.ResponseWriter, r *http.Request)  {

	start := "{"
	content := ""
	end := "}"

	file, _, err := r.FormFile("file")
	if err != nil {
		content = "\"error\":1,\"result\":{\"msg\":\"指定了无效的文件\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		content = "\"error\":1,\"result\":{\"msg\":\"无法读取文件内容\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}

	filetype := http.DetectContentType(fileBytes)
	//log.Println("filetype = " + filetype)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
		break
	default:
		content = "\"error\":1,\"result\":{\"msg\":\"文件类型错误\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}

	fileName := randToken(12)	// 指定文件名
	fileEndings, err := mime.ExtensionsByType(filetype)	// 获取文件扩展名
	//log.Println("fileEndings = " + fileEndings[0])
	// 指定文件存储路径
	newPath := filepath.Join("web", "static", "photo", fileName + fileEndings[0])
	//fmt.Printf("FileType: %s, File: %s\n", filetype, newPath)

	newFile, err := os.Create(newPath)
	if err != nil {
		log.Println("创建文件失败：" + err.Error())
		content = "\"error\":1,\"result\":{\"msg\":\"创建文件失败\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		log.Println("写入文件失败：" + err.Error())
		content = "\"error\":1,\"result\":{\"msg\":\"保存文件内容失败\",\"path\":\"\"}"
		w.Write([]byte(start + content + end))
		return
	}

	path := "/static/photo/" + fileName + fileEndings[0]
	content = "\"error\":0,\"result\":{\"fileType\":\"image/png\",\"path\":\"" + path + "\",\"fileName\":\"ce73ac68d0d93de80d925b5a.png\"}"
	w.Write([]byte(start + content + end))
	return
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

