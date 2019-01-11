package audience

import (
	"bytes"
	"encoding/json"
	"errors"
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
	url     string
	confirm string
}

// New returns new instance of uploader.
func New() *Uploader {
	return &Uploader{
		url:     "https://api-audience.yandex.ru/v1/management/segments/upload_csv_file",
		confirm: "https://api-audience.yandex.ru/v1/management/segment/%d/confirm",
	}
}

func (u *Uploader) confirmURL(id int) string {
	return fmt.Sprintf(u.confirm, id)
}

// Do processes file upload.
func (u *Uploader) Do(token, name, seg string, file multipart.File) (*Payload, error) {
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read content: %v", err)
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// fileWriter, err := newpart(writer, name, "text/csv")
	// if err != nil {
	// 	return nil, fmt.Errorf("new part: %v", err)
	// }
	// io.Copy(fileWriter, file)

	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return nil, fmt.Errorf("new part: %v", err)
	}

	// _, err = fileWriter.Write(content)
	_, err = part.Write(content)
	if err != nil {
		return nil, fmt.Errorf("file write: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("writer close: %v", err)
	}

	request, err := http.NewRequest("POST", u.url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", u.url)
	}
	request.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))
	request.Header.Add("Content-Type", writer.FormDataContentType())

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		return nil, fmt.Errorf("dump request: %v", err)
	}
	fmt.Println(string(dump))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(e)
		if err != nil {
			return nil, fmt.Errorf("decode error: %v", err)
		}
		return nil, errors.New(e.Message)
	}

	p := &Payload{}
	err = json.NewDecoder(resp.Body).Decode(p)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %v", err)
	}
	fmt.Println(p.Print())

	s := &Confirm{
		Segment: SegmentConfirm{
			ID:          p.Segment.ID,
			Name:        seg,
			Hashed:      0,
			ContentType: "mac",
		},
	}

	err = json.NewEncoder(body).Encode(s)
	if err != nil {
		return nil, fmt.Errorf("encode confirm: %v", err)
	}

	url := u.confirmURL(s.Segment.ID)
	request, err = http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", url)
	}
	request.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))

	dump, err = httputil.DumpRequest(request, true)
	if err != nil {
		return nil, fmt.Errorf("dump request: %v", err)
	}
	fmt.Println(string(dump))

	var respC *http.Response
	respC, err = client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send request: %v", err)
	}
	defer respC.Body.Close()

	if respC.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		err = json.NewDecoder(respC.Body).Decode(e)
		if err != nil {
			return nil, fmt.Errorf("decode error: %v", err)
		}
		return nil, errors.New(e.Message)
	}

	err = json.NewDecoder(respC.Body).Decode(p)
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
