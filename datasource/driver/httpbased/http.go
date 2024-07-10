package httpbased

import (
	"github.com/allape/goview/datasource/driver"
	"io"
	"net/http"
)

// region options

type StopSign bool

type Option interface {
	Apply(req *http.Request, res *http.Response) StopSign
}

type RedirectWriterOption struct {
	Writer io.Writer
	Option
}

func (o RedirectWriterOption) Apply(_ *http.Request, res *http.Response) StopSign {
	_, _ = io.Copy(o.Writer, res.Body)
	return true
}

type HeaderOnlyOption struct {
	Option
}

func (o HeaderOnlyOption) Apply(_ *http.Request, _ *http.Response) StopSign {
	return true
}

// endregion

type AbstractDriver struct {
	driver.Driver

	client *http.Client
}

func (d *AbstractDriver) NewRequest(url string, options ...Option) ([]byte, http.Header, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Accept", "application/json; charset=utf-8")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	for _, option := range options {
		if option.Apply(req, resp) {
			return nil, resp.Header, nil
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.Header, err
	}

	return body, resp.Header, nil
}

func (d *AbstractDriver) Setup(client *http.Client) error {
	if client == nil {
		client = &http.Client{}
	}
	d.client = client
	return nil
}

//func (d *HttpBasedDriver) List(wd string) ([]File, error) {
//
//}
//
//func (d *HttpBasedDriver) Concatenate(file string, writer io.Writer) error {
//
//}
