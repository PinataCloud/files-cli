package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func Upload(filePath string, groupId string, name string, verbose bool) (UploadResponse, error) {
	jwt, err := findToken()
	if err != nil {
		return UploadResponse{}, err
	}

	stats, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("File or folder does not exist")
		return UploadResponse{}, errors.Join(err, errors.New("file or folder does not exist"))
	}

	files, err := pathsFinder(filePath, stats)
	if err != nil {
		return UploadResponse{}, err
	}

	body := &bytes.Buffer{}
	contentType, err := createMultipartRequest(filePath, files, body, stats, groupId, name)
	if err != nil {
		return UploadResponse{}, err
	}

	var requestBody io.Reader
	if !verbose {
		requestBody = body
	} else {
		totalSize := int64(body.Len())
		fmt.Printf("Uploading %s (%s)\n", stats.Name(), formatSize(int(totalSize)))
		requestBody = newProgressReader(body, totalSize)
	}

	url := fmt.Sprintf("https://uploads.pinata.cloud/v3/files")
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return UploadResponse{}, errors.Join(err, errors.New("failed to create the request"))
	}
	req.Header.Set("Authorization", "Bearer "+string(jwt))
	req.Header.Set("content-type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UploadResponse{}, errors.Join(err, errors.New("failed to send the request"))
	}
	if resp.StatusCode != 200 {
		return UploadResponse{}, fmt.Errorf("server Returned an error %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var response UploadResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return UploadResponse{}, err
	}

	formattedJSON, err := json.MarshalIndent(response.Data, "", "    ")
	if err != nil {
		return UploadResponse{}, errors.New("failed to format JSON")
	}

	fmt.Println(string(formattedJSON))

	return response, nil
}

type progressReader struct {
	r   io.Reader
	bar *progressbar.ProgressBar
}

func cmpl() {
	fmt.Println()
}

func newProgressReader(r io.Reader, size int64) *progressReader {
	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription("Uploading..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
		progressbar.OptionOnCompletion(cmpl),
	)
	return &progressReader{r: r, bar: bar}
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)
	if err != nil {
		return 0, err
	}
	err = pr.bar.Add(n)
	if err != nil {
		return 0, err
	}
	return
}

func formatSize(bytes int) string {
	const (
		KB = 1000
		MB = KB * KB
		GB = MB * KB
	)

	var formattedSize string

	switch {
	case bytes < KB:
		formattedSize = fmt.Sprintf("%d bytes", bytes)
	case bytes < MB:
		formattedSize = fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	case bytes < GB:
		formattedSize = fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	default:
		formattedSize = fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	}

	return formattedSize
}

func createMultipartRequest(filePath string, files []string, body io.Writer, stats os.FileInfo, groupId string, name string) (string, error) {
	contentType := ""
	writer := multipart.NewWriter(body)

	fileIsASingleFile := !stats.IsDir()
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return contentType, err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal("could not close file")
			}
		}(file)

		var part io.Writer
		if fileIsASingleFile {
			part, err = writer.CreateFormFile("file", filepath.Base(f))
		} else {
			relPath, _ := filepath.Rel(filePath, f)
			part, err = writer.CreateFormFile("file", filepath.Join(stats.Name(), relPath))
		}
		if err != nil {
			return contentType, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return contentType, err
		}
	}

	if groupId != "" {
		err := writer.WriteField("group_id", groupId)
		if err != nil {
			return contentType, err
		}
	}

	nameToUse := stats.Name()
	if name != "nil" {
		nameToUse = name
	}
	err := writer.WriteField("name", nameToUse)
	if err != nil {
		return contentType, err
	}

	err = writer.Close()
	if err != nil {
		return contentType, err
	}

	contentType = writer.FormDataContentType()

	return contentType, nil
}

func pathsFinder(filePath string, stats os.FileInfo) ([]string, error) {
	var err error
	files := make([]string, 0)
	fileIsASingleFile := !stats.IsDir()
	if fileIsASingleFile {
		files = append(files, filePath)
		return files, err
	}
	err = filepath.Walk(filePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return files, err
}
