package fits

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"sort"
)

type Method struct {
	MethodId    string `json:"methodID"`    // A valid method identifier for observation type e.g., doas-s.
	Name        string `json:"name"`        // The method name e.g., Bernese v5.0.
	Description string `json:"description"` // A description of the method e.g., Bernese v5.0 GNS processing software.
	Reference   string `json:"reference"`   // A link to further information about the method.
}

type Methods []Method

func (m Methods) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return m.Write(f)
}

func (m Methods) Write(wr io.Writer) error {

	lines := [][]string{
		{"Id", "Name", "Description", "Reference"},
	}
	for _, v := range m {
		lines = append(lines, []string{
			v.MethodId, v.Name, v.Description, v.Reference,
		})
	}

	return csv.NewWriter(wr).WriteAll(lines)
}

func (f Fits) Methods(ctx context.Context, extra ...string) ([]Method, error) {

	v := url.Values{}
	if len(extra) > 0 {
		v.Set("typeID", extra[0])
	}

	data, err := f.Query(ctx, "/method", v, "application/json;version=1")
	if err != nil {
		return nil, err
	}

	var types Results
	if err := json.Unmarshal(data, &types); err != nil {
		return nil, err
	}

	list := make(map[Method]interface{})
	for _, m := range types.Methods {
		list[m] = true
	}

	var res []Method
	for k := range list {
		res = append(res, k)
	}

	sort.Slice(res, func(i, j int) bool {
		switch {
		case res[i].Name < res[j].Name:
			return true
		case res[i].Name > res[j].Name:
			return false
		default:
			return res[i].MethodId < res[j].MethodId
		}
	})

	return res, nil
}
