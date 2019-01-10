package audience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
)

// Uploader does csv upload.
type Uploader struct {
	url string
}

// New returns new instance of uploader.
func New() *Uploader {
	return &Uploader{
		url: "https://api-audience.yandex.ru/v1/management/segments/upload_csv_file",
	}
}

// Do processes file upload.
func (u *Uploader) Do(token, name string, file multipart.File) (*Payload, error) {
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read content: %v", err)
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	fileWriter, err := newpart(writer, name, "text/csv")
	if err != nil {
		return nil, fmt.Errorf("new part: %v", err)
	}
	io.Copy(fileWriter, file)

	_, err = fileWriter.Write(content)
	if err != nil {
		return nil, fmt.Errorf("file write: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("writer close: %v", err)
	}

	request, err := http.NewRequest("POST", u.url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))
	request.Header.Add("Content-Type", writer.FormDataContentType())

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(dump))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wront status code: %d", resp.StatusCode)
	}

	p := &Payload{}
	err = json.NewDecoder(resp.Body).Decode(p)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %v", err)
	}
	fmt.Println(p.Print())

	return p, nil
}

func newpart(w *multipart.Writer, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
