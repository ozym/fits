package fits

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Observation struct {
	Method string
	Type   string
	Site   string

	Time  time.Time // The date-time of the observation in ISO8601 format, UTC time zone.
	Value float64   // The observation value.
	Error float64   // The observation error. 0 is used for an unknown error.
}

type Observations []Observation

func (o Observations) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return o.Write(f)
}

func (o Observations) Write(wr io.Writer) error {

	lines := [][]string{
		{"Method", "Type", "Site", "Time", "Vaue", "Error"},
	}
	for _, v := range o {
		lines = append(lines, []string{
			v.Method, v.Type, v.Site, v.Time.Format(time.RFC3339), toStr(v.Value), toStr(v.Error),
		})
	}

	return csv.NewWriter(wr).WriteAll(lines)
}

func (f Fits) Observations(ctx context.Context, methodId, typeId, siteId string) ([]Observation, error) {

	v := url.Values{}
	v.Set("typeID", typeId)
	v.Set("methodID", methodId)
	v.Set("siteID", siteId)
	if f.Days > 0 {
		v.Set("days", strconv.Itoa(f.Days))
	}

	data, err := f.Query(ctx, "/observation", v, "text/csv;version=1")
	if err != nil {
		return nil, err
	}

	lines, err := csv.NewReader(bytes.NewBuffer(data)).ReadAll()
	if err != nil {
		return nil, err
	}

	if !(len(lines) > 1) {
		return nil, nil
	}

	var obs []Observation
	for _, l := range lines[1:] {
		if len(l) < 3 {
			continue
		}

		t, err := time.Parse(time.RFC3339, l[0])
		if err != nil {
			return nil, err
		}

		v, err := strconv.ParseFloat(l[1], 64)
		if err != nil {
			return nil, err
		}

		e, err := strconv.ParseFloat(l[2], 64)
		if err != nil {
			return nil, err
		}

		obs = append(obs, Observation{
			Type:   typeId,
			Method: methodId,
			Site:   siteId,
			Time:   t,
			Value:  v,
			Error:  e,
		})
	}

	return obs, nil
}
