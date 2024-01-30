package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 解析请求体中的JSON数据
type PostData struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Lrc  string `json:"lrc"`
}

var folder = "music/"

func downloadFile(url, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}

	return nil
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 允许跨域请求
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		// 处理预检请求
		return
	}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var data PostData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	url := data.URL

	filename := folder + data.Name // 可以根据需求修改保存的文件名

	log.Println(data.Name + " download...")

	if url == "" {
		http.Error(w, "get url error", http.StatusBadRequest)
		return
	}
	writeToFile(filename+".lrc", data.Lrc)
	err = downloadFile(url, filename+".mp3")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to download file: %s", err), http.StatusInternalServerError)
		return
	}

	log.Println(data.Name + " downloaded successfully")
	fmt.Fprintf(w, "File downloaded successfully")
}

func main() {
	http.HandleFunc("/download", downloadHandler)
	log.Println("http://localhost:8182")
	http.ListenAndServe(":8182", nil)
}
