package fits

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
)

type Site struct {
	SiteId             string  `json:"siteID"` // Site identifier e.g., WI000.
	Name               string  `json:"name"`   // Site name e.g, White Island Volcano.
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Height             float64 `json:"height"`             // Site height (m).
	GroundRelationship float64 `json:"groundRelationship"` // Site ground relationship (m). Sites above ground level have a negative ground relationship.
}

type Sites []Site

func (s Sites) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return s.Write(f)
}

func (s Sites) Write(wr io.Writer) error {

	lines := [][]string{
		{"Id", "Name", "Latitude", "Longitude", "Height", "Ground Relationship"},
	}
	for _, v := range s {
		lines = append(lines, []string{
			v.SiteId, v.Name, toStr(v.Latitude), toStr(v.Longitude), toStr(v.Height), toStr(v.GroundRelationship),
		})
	}

	return csv.NewWriter(wr).WriteAll(lines)
}

type Geometry struct {
	Coordinates [2]float64 `json:"coordinates"`
	Type        string     `json:"type"`
}

type Feature struct {
	Geom       Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
	Type       string                 `json:"type"`
}

type Features struct {
	Features []Feature `json:"features"`
	Type     string    `json:"type"`
}

func (f Fits) Sites(ctx context.Context, extra ...string) ([]Site, error) {

	v := url.Values{}
	if len(extra) > 0 {
		v.Set("typeID", extra[0])
	}

	data, err := f.Query(ctx, "/site", v, "application/vnd.geo+json;version=1")
	if err != nil {
		return nil, err
	}

	var features Features
	if err := json.Unmarshal(data, &features); err != nil {
		return nil, err
	}
	if features.Type != "FeatureCollection" {
		return nil, fmt.Errorf("invalid features type: %s", features.Type)
	}

	var sites []Site
	for _, f := range features.Features {
		if f.Type != "Feature" {
			return nil, fmt.Errorf("invalid site type: %s", f.Type)
		}
		if f.Geom.Type != "Point" {
			return nil, fmt.Errorf("invalid site geom type: %s", f.Geom.Type)
		}
		s, ok := f.Properties["siteID"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid site id: %v", f.Properties["siteID"])
		}

		n, ok := f.Properties["name"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid site name: %v", f.Properties["name"])
		}

		h, ok := f.Properties["height"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid site height: %v", f.Properties["height"])
		}

		g, ok := f.Properties["groundRelationship"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid site ground relationship: %v", f.Properties["groundRelationship"])
		}

		sites = append(sites, Site{
			SiteId:             s,
			Name:               n,
			Latitude:           f.Geom.Coordinates[1],
			Longitude:          f.Geom.Coordinates[0],
			Height:             h,
			GroundRelationship: g,
		})
	}

	sort.Slice(sites, func(i, j int) bool { return sites[i].SiteId < sites[j].SiteId })

	return sites, nil
}
