package domain

type CatalogEntry struct {
	Code, Label, Description string
	Active                   bool
	Rank                     int
}
type Catalog struct{ Regions, Venues, Themes []CatalogEntry }

func DefaultCatalog() Catalog {
	return Catalog{Regions: []CatalogEntry{{"east", "East", "Eastern region", true, 1}, {"west", "West", "Western region", true, 2}, {"south", "South", "Southern region", true, 3}}, Venues: []CatalogEntry{{"hall", "Hall", "Indoor hall", true, 1}, {"garden", "Garden", "Outdoor garden", true, 2}}, Themes: []CatalogEntry{{"classic", "Classic", "Traditional ceremony", true, 1}, {"modern", "Modern", "Contemporary ceremony", true, 2}}}
}
func (c Catalog) Region(code string) (CatalogEntry, bool) {
	for _, x := range c.Regions {
		if x.Code == code && x.Active {
			return x, true
		}
	}
	return CatalogEntry{}, false
}
func (c Catalog) Venue(code string) (CatalogEntry, bool) {
	for _, x := range c.Venues {
		if x.Code == code && x.Active {
			return x, true
		}
	}
	return CatalogEntry{}, false
}
func (c Catalog) Theme(code string) (CatalogEntry, bool) {
	for _, x := range c.Themes {
		if x.Code == code && x.Active {
			return x, true
		}
	}
	return CatalogEntry{}, false
}
func (c Catalog) AllActive() []CatalogEntry {
	out := []CatalogEntry{}
	for _, list := range [][]CatalogEntry{c.Regions, c.Venues, c.Themes} {
		for _, x := range list {
			if x.Active {
				out = append(out, x)
			}
		}
	}
	return out
}
func (c Catalog) ValidRecord(r Record) bool {
	_, ok := c.Region(r.Region)
	if r.Region != "" && !ok {
		return false
	}
	if r.Venue != "" {
		if _, ok = c.Venue(r.Venue); !ok {
			return false
		}
	}
	return r.Validate() == nil
}
