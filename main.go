package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"os"
	"time"

	"github.com/danikarik/yandex-audience/pkg/web"

	"github.com/danikarik/yandex-audience/pkg/config"

	"github.com/gorilla/sessions"
)

// type config struct {
// 	Token string
// }

// Payload is a response json.
type Payload struct {
	Segment Segment `json:"segment"`
}

// Print prints payload.
func (p *Payload) Print() string {
	format := `
id: %d,
type: %s,
status: %s,
has_guests: %t,
guest_quantity: %d,
can_create_dependent: %t,
has_derivatives: %t,
hashed: %t,
item_quantity: %d,
guest: %t
`
	return fmt.Sprintf(format,
		p.Segment.ID,
		p.Segment.Type,
		p.Segment.Status,
		p.Segment.HasGuests,
		p.Segment.GuestQuantity,
		p.Segment.CanCreateDependent,
		p.Segment.HasDerivatives,
		p.Segment.Hashed,
		p.Segment.ItemQuantity,
		p.Segment.Guest,
	)
}

// Segment represents segment params.
type Segment struct {
	ID                 int    `json:"id"`
	Type               string `json:"type"`
	Status             string `json:"status"`
	HasGuests          bool   `json:"has_guests"`
	GuestQuantity      int    `json:"guest_quantity"`
	CanCreateDependent bool   `json:"can_create_dependent"`
	HasDerivatives     bool   `json:"has_derivatives"`
	Hashed             bool   `json:"hashed"`
	ItemQuantity       int    `json:"item_quantity"`
	Guest              bool   `json:"guest"`
}

var (
	// ErrTokenNotFound will be returned if token is not found.
	ErrTokenNotFound = errors.New("could not find YANDEX_OAUTH_TOKEN environment variable")
	// ErrBadRequest will be returned when response status code is 400.
	ErrBadRequest = errors.New("400: request failed")
)

func main() {
	// Load environment variables with no prefix.
	spec := loadenv("")
	// Create session store.
	store := sessions.NewCookieStore([]byte(spec.SessionKey))
	// Create new web container.
	gob.Register(time.Now())
	web := web.New(store, spec)
	if err := http.ListenAndServe(spec.Addr(), web.Handler()); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
		os.Exit(0)
	}

	// conf, err := loadenv()
	// if err != nil {
	// 	log.Fatalf("loadenv: %v", err)
	// }
	// p := "./data/data.csv"
	// err = upload("https://api-audience.yandex.ru/v1/management/segments/upload_csv_file", conf.Token, p)
	// if err != nil {
	// 	log.Fatalf("upload: %v", err)
	// }
}

func loadenv(prefix string) config.Specification {
	spec, err := config.NewSpec(prefix)
	if err != nil {
		uerr := config.Usage(prefix, spec)
		if uerr != nil {
			log.Fatalf("could not load environment variables: %v", uerr)
		}
		os.Exit(0)
	}
	return spec
}

func upload(url, token, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	fileWriter, err := newpart(writer, stat.Name(), "text/csv")
	if err != nil {
		return err
	}
	io.Copy(fileWriter, file)

	_, err = fileWriter.Write(content)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token))
	request.Header.Add("Content-Type", writer.FormDataContentType())

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		return err
	}
	fmt.Println(string(dump))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrBadRequest
	}

	p := &Payload{}
	err = json.NewDecoder(resp.Body).Decode(p)
	if err != nil {
		return err
	}
	fmt.Println(p.Print())

	return nil
}

func newpart(w *multipart.Writer, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}
