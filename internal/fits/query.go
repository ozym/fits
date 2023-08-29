package fits

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Fits struct {
	Service string
	Scheme  string
	Days    int
}

func Scheme(scheme string) func(*Fits) {
	return func(f *Fits) {
		f.Scheme = scheme
	}
}

func Days(days int) func(*Fits) {
	return func(f *Fits) {
		f.Days = days
	}
}

// NewFits returns a default Event with non-standard values set
// by using the option parameter functions.
func New(service string, opts ...func(*Fits)) Fits {
	client := Fits{
		Service: service,
		Scheme:  "https",
	}

	for _, opt := range opts {
		opt(&client)
	}

	return client
}

func (f *Fits) Set(opts ...func(*Fits)) {
	for _, opt := range opts {
		opt(f)
	}
	return
}

func (f Fits) Query(ctx context.Context, path string, values url.Values, accept string) ([]byte, error) {
	u := url.URL{
		Scheme:   f.Scheme,
		Host:     f.Service,
		Path:     path,
		RawQuery: values.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf(string(body))
	}

	return body, nil
}
