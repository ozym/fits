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

// Type holds an observation type details.
type Type struct {
	TypeId      string `json:"typeID"`      // A type identifier for observations e.g., e.
	Name        string `json:"name"`        //Type name e.g., east
	Description string `json:"description"` // Type description e.g., displacement from initial position
	Unit        string `json:"unit"`        // Type unit e.g., mm.
}

type Types []Type

func (t Types) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Write(f)
}

func (t Types) Write(wr io.Writer) error {

	lines := [][]string{
		{"Id", "Name", "Description", "Unit"},
	}
	for _, v := range t {
		lines = append(lines, []string{
			v.TypeId, v.Name, v.Description, v.Unit,
		})
	}

	return csv.NewWriter(wr).WriteAll(lines)
}

func (f Fits) Types(ctx context.Context) ([]Type, error) {

	data, err := f.Query(ctx, "/type", url.Values{}, "application/json;version=1")
	if err != nil {
		return nil, err
	}

	var types Results
	if err := json.Unmarshal(data, &types); err != nil {
		return nil, err
	}

	list := make(map[Type]interface{})
	for _, t := range types.Types {
		list[t] = true
	}

	var res []Type
	for k := range list {
		res = append(res, k)
	}

	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })

	return res, nil
}
