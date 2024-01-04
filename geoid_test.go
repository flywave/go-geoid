package geoid

import "testing"

func TestGeoid(t *testing.T) {

	lat := float64(36.35697678803291)
	lon := float64(112.84426434425896)
	alt := float64(347)

	g := NewGeoid(EGM2008, false)
	galt := g.ConvertHeight(lat, lon, alt, GEOIDTOELLIPSOID)

	if galt == 0 {
		t.FailNow()
	}
}
