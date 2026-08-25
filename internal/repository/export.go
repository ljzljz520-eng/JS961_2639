package repository

import (
	"encoding/json"
	"sort"
	"weddingtemplates/internal/domain"
)

func (r *Repository) Export(q domain.Query) ([]byte, error) {
	p, e := r.Search(q)
	if e != nil {
		return nil, e
	}
	sort.Slice(p.Items, func(i, j int) bool { return p.Items[i].ID < p.Items[j].ID })
	return json.Marshal(p.Items)
}
func (r *Repository) Metrics() (domain.Metrics, error) {
	items, e := r.Store.ListRecords()
	if e != nil {
		return domain.Metrics{}, e
	}
	return domain.ComputeMetrics(items), nil
}
func (r *Repository) IDs() ([]string, error) {
	items, e := r.Store.ListRecords()
	if e != nil {
		return nil, e
	}
	out := make([]string, 0, len(items))
	for _, x := range items {
		out = append(out, x.ID)
	}
	sort.Strings(out)
	return out, nil
}
func (r *Repository) Remove(id string) error { return r.Store.DeleteRecord(id) }
