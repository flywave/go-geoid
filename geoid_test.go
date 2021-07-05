package geoid

import "testing"

func TestGeoid(t *testing.T) {
	SetGeoidPath("./data")

	lat := float64(10)
	lon := float64(20)
	alt := float64(30)

	g := NewGeoid(EGM84, false)
	galt := g.ConvertHeight(lat, lon, alt, ELLIPSOIDTOGEOID)

	if galt == 0 {
		t.FailNow()
	}
}
