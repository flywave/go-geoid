#include "geoid_api.h"

#include <GeographicLib/Geoid.hpp>

#include <algorithm>
#include <fstream>
#include <memory>
#include <stdint.h>
#include <vector>

#ifdef __cplusplus
extern "C" {
#endif

static std::string geoid_path;

void set_geoid_path(const char *path) { geoid_path = path; }

static std::unique_ptr<GeographicLib::Geoid> get_geoid(vertical_datum g,
                                                       bool cubic) {
  switch (g) {
  case vertical_datum::EGM84:
    return std::make_unique<GeographicLib::Geoid>("egm84-30", geoid_path, cubic,
                                                  true);
  case vertical_datum::EGM2008:
    return std::make_unique<GeographicLib::Geoid>("egm2008-5", geoid_path,
                                                  cubic, true);
  default:
    return std::make_unique<GeographicLib::Geoid>("egm96-15", geoid_path, cubic,
                                                  true);
  }
}

struct _geoid_t {
  std::unique_ptr<GeographicLib::Geoid> id;
};

geoid_t *geoid_new(int g, _Bool cubic) {
  return new geoid_t{get_geoid(static_cast<vertical_datum>(g), cubic)};
}

void geoid_free(geoid_t *gid) { delete gid; }

void geoid_cache_area(geoid_t *gid, double south, double west, double north,
                      double east) {
  gid->id->CacheArea(south, west, north, east);
}

void geoid_cache_all(geoid_t *gid) { gid->id->CacheAll(); }

void geoid_cache_clear(geoid_t *gid) { gid->id->CacheClear(); }

double geoid_get_height(geoid_t *gid, double lat, double lon) {
  return (*gid->id)(lat, lon);
}

double geoid_convert_height(geoid_t *gid, double lat, double lon, double h,
                            int flag) {
  gid->id->ConvertHeight(lat, lon, h,
                         static_cast<GeographicLib::Geoid::convertflag>(flag));
}

_Bool geoid_is_cached(geoid_t *gid) { return gid->id->Cache(); }

void geoid_get_cache_area(geoid_t *gid, double *south, double *west,
                          double *north, double *east) {
  *south = gid->id->CacheSouth();
  *west = gid->id->CacheWest();
  *north = gid->id->CacheNorth();
  *east = gid->id->CacheEast();
}

const char *geoid_get_name(geoid_t *gid) {
  return gid->id->GeoidName().c_str();
}

double geoid_get_major_radius(geoid_t *gid) { return gid->id->MajorRadius(); }

double geoid_get_flattening(geoid_t *gid) { return gid->id->Flattening(); }

#ifdef __cplusplus
}
#endif
