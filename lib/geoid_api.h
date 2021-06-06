#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

#if defined(WIN32) || defined(WINDOWS) || defined(_WIN32) || defined(_WINDOWS)
#define FLYWAVE_GEOID_API __declspec(dllexport)
#else
#define FLYWAVE_GEOID_API
#endif

enum vertical_datum { HAE, EGM84, EGM96, EGM2008, UNKNOWN };

enum convertflag {
  ELLIPSOIDTOGEOID = -1,
  NONE = 0,
  GEOIDTOELLIPSOID = 1,
};

typedef struct _geoid_t geoid_t;

FLYWAVE_GEOID_API void set_geoid_path(const char *path);

FLYWAVE_GEOID_API geoid_t *geoid_new(int g, _Bool cubic);

FLYWAVE_GEOID_API void geoid_free(geoid_t *gid);

FLYWAVE_GEOID_API void geoid_cache_area(geoid_t *gid, double south, double west,
                                        double north, double east);

FLYWAVE_GEOID_API void geoid_cache_all(geoid_t *gid);

FLYWAVE_GEOID_API void geoid_cache_clear(geoid_t *gid);

FLYWAVE_GEOID_API double geoid_get_height(geoid_t *gid, double lat, double lon);

FLYWAVE_GEOID_API double geoid_convert_height(geoid_t *gid, double lat,
                                              double lon, double h, int flag);

FLYWAVE_GEOID_API _Bool geoid_is_cached(geoid_t *gid);

FLYWAVE_GEOID_API void geoid_get_cache_area(geoid_t *gid, double *south,
                                            double *west, double *north,
                                            double *east);

FLYWAVE_GEOID_API const char *geoid_get_name(geoid_t *gid);

FLYWAVE_GEOID_API double geoid_get_major_radius(geoid_t *gid);

FLYWAVE_GEOID_API double geoid_get_flattening(geoid_t *gid);

#ifdef __cplusplus
}
#endif
