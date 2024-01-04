package geoid

// #include <stdlib.h>
// #include <string.h>
// #include "geoid_api.h"
// #cgo CFLAGS:   -I ./lib
// #cgo CXXFLAGS:   -I ./lib
// #cgo linux LDFLAGS:  -L ./lib/linux -Wl,--start-group  -lstdc++ -lm -pthread -ldl -lcgeoid -lgeographic -Wl,--end-group
// #cgo windows LDFLAGS:  -L ./lib/windows -Wl,--start-group  -lstdc++ -lm   -lcgeoid -lgeographic -Wl,--end-group
// #cgo darwin LDFLAGS: -L　./lib/darwin -lcgeoid -lgeographic
// #cgo darwin ,arm LDFLAGS: -L　./lib/darwin_arm -lcgeoid -lgeographic
import "C"
import (
	"path/filepath"
	"runtime"
	"unsafe"
)

type Geoid struct {
	m *C.struct__geoid_t
}

func getCurrentDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

func init() {
	dir := getCurrentDir()
	SetGeoidPath(filepath.Join(dir, "./data"))
}

func SetGeoidPath(path string) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	C.set_geoid_path(cpath)
}

func NewGeoid(g VerticalDatum, cubic bool) *Geoid {
	ret := &Geoid{m: C.geoid_new(C.int(g), C.bool(cubic))}
	runtime.SetFinalizer(ret, (*Geoid).free)
	return ret
}

func (g *Geoid) free() {
	C.geoid_free(g.m)
}

func (g *Geoid) GetName() string {
	return C.GoString(C.geoid_get_name(g.m))
}

func (g *Geoid) GetMajorRadius() float64 {
	return float64(C.geoid_get_major_radius(g.m))
}

func (g *Geoid) GetFlattening() float64 {
	return float64(C.geoid_get_flattening(g.m))
}

func (g *Geoid) CacheArea(south, west, north, east float64) {
	C.geoid_cache_area(g.m, C.double(south), C.double(west), C.double(north), C.double(east))
}

func (g *Geoid) CacheAll() {
	C.geoid_cache_all(g.m)
}

func (g *Geoid) CacheClear() {
	C.geoid_cache_clear(g.m)
}

func (g *Geoid) IsCached() bool {
	return bool(C.geoid_is_cached(g.m))
}

func (g *Geoid) GetCacheArea() (south, west, north, east float64) {
	var csouth, cwest, cnorth, ceast C.double
	C.geoid_get_cache_area(g.m, &csouth, &cwest, &cnorth, &ceast)
	south, west, north, east = float64(csouth), float64(cwest), float64(cnorth), float64(ceast)
	return
}

func (g *Geoid) GetHeight(lat, lon float64) float64 {
	return float64(C.geoid_get_height(g.m, C.double(lat), C.double(lon)))
}

func (g *Geoid) ConvertHeight(lat, lon, height float64, flag ConvertFlag) float64 {
	return float64(C.geoid_convert_height(g.m, C.double(lat), C.double(lon), C.double(height), C.int(flag)))
}
