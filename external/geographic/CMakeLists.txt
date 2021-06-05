project (GeographicLib)

# Version information
set (PROJECT_VERSION_MAJOR 1)
set (PROJECT_VERSION_MINOR 44)
set (PROJECT_VERSION_PATCH 0)
set (PROJECT_VERSION "${PROJECT_VERSION_MAJOR}.${PROJECT_VERSION_MINOR}")
if (PROJECT_VERSION_PATCH GREATER 0)
  set (PROJECT_VERSION "${PROJECT_VERSION}.${PROJECT_VERSION_PATCH}")
endif ()

set (LIBVERSION_API 14)
set (LIBVERSION_BUILD 14.2.0)
string (TOLOWER ${PROJECT_NAME} PROJECT_NAME_LOWER)
string (TOUPPER ${PROJECT_NAME} PROJECT_NAME_UPPER)

cmake_minimum_required (VERSION 2.8.4) # This version was released 2011-02-16

# User-settable variables

# (1) COMMON_INSTALL_PATH governs the installation convention.  If it
# is on ON (the Linux default), the installation is to a common
# directory, e.g., /usr/local.  If it is OFF (the Windows default),
# the installation directory contains the package name, e.g.,
# c:/pkg/GeographicLib-1.22.  The installation directories for the
# documentation, cmake configuration, python and matlab interfaces all
# depend on the variable with deeper paths relative to
# CMAKE_INSTALL_PREFIX being used when it's ON.

if (WIN32)
  option (COMMON_INSTALL_PATH "Use a common installation path for packages" OFF)
else ()
  option (COMMON_INSTALL_PATH "Use a common installation path for packages" ON)
endif ()

# The use of PACKAGE_PATH and INSTALL_PATH is now DEPRECATED.
# (2) PACKAGE_PATH and INSTALL_PATH govern the find_package search
# path and the installation directory.  (find_package is not used by
# GeographicLib since it doesn't depend on other packages.  However
# PACKAGE_PATH is used here for uniformity with other packages which
# adopt the same conventions.)
#
# If PACKAGE_PATH is defined, it is prepended to CMAKE_PREFIX_PATH.
#
# If INSTALL_PATH is not specified but PACKAGE_PATH is, then
# INSTALL_PATH is set to
#   ${PACKAGE_PATH}, if COMMON_INSTALL_PATH is ON;
#   ${PACKAGE_PATH}/${PROJECT_NAME}-${PROJECT_VERSION}, otherwise.
#
# If INSTALL_PATH is now defined, then set CMAKE_INSTALL_PREFIX to
# INSTALL_PATH.
#
# Typically, only PACKAGE_PATH needs to be specified, e.g.,
# cmake -D PACKAGE_PATH=/opt .. (on Linux)
#   => CMAKE_PREFIX_PATH=/opt   CMAKE_INSTALL_PREFIX=/opt
# cmake -D PACKAGE_PATH=C:/pkg .. (on Windows)
#   => CMAKE_PREFIX_PATH=C:/pkg CMAKE_INSTALL_PREFIX=C:/pkg/GeographicLib-1.22

if (PACKAGE_PATH)
  set (CMAKE_PREFIX_PATH ${PACKAGE_PATH} ${CMAKE_PREFIX_PATH})
  message (STATUS "CMAKE_PREFIX_PATH set to ${CMAKE_PREFIX_PATH}")
endif ()

if (NOT INSTALL_PATH AND PACKAGE_PATH)
  if (COMMON_INSTALL_PATH)
    set (INSTALL_PATH ${PACKAGE_PATH} CACHE PATH "Installation directory" FORCE)
  else ()
    set (INSTALL_PATH ${PACKAGE_PATH}/${PROJECT_NAME}-${PROJECT_VERSION}
      CACHE PATH "Installation directory" FORCE)
  endif ()
endif ()
if (INSTALL_PATH)
  file (TO_CMAKE_PATH ${INSTALL_PATH} CMAKE_INSTALL_PREFIX)
  message (STATUS "CMAKE_INSTALL_PREFIX set to ${CMAKE_INSTALL_PREFIX}")
endif ()

# (3) Where to look for data files.  Various classes look in the geoids,
# gravity, magnetic, subdirectories of ${GEOGRAPHICLIB_DATA}.
if (WIN32)
  # The binary installers for the data files for Windows are created
  # with Inno Setup which uses {commonappdata} which (since Windows
  # Vista) is C:/ProgramData.
  set (GEOGRAPHICLIB_DATA
    "C:/ProgramData/GeographicLib"
    CACHE PATH "Location for data for GeographicLib")
else ()
  set (GEOGRAPHICLIB_DATA
    "/usr/local/share/GeographicLib"
    CACHE PATH "Location for data for GeographicLib")
endif ()

# (4) Build which libraries?  Possible values are SHARED, STATIC, BOTH.
if (MSVC)
  set (GEOGRAPHICLIB_LIB_TYPE STATIC CACHE STRING
    "Types of library generated: SHARED, STATIC (default), or BOTH")
else ()
  set (GEOGRAPHICLIB_LIB_TYPE SHARED CACHE STRING
    "Types of library generated: SHARED (default), STATIC, or BOTH")
endif ()
set_property (CACHE GEOGRAPHICLIB_LIB_TYPE
  PROPERTY STRINGS "SHARED" "STATIC" "BOTH")

set (GEOGRAPHICLIB_SHARED_LIB OFF)
set (GEOGRAPHICLIB_STATIC_LIB ON)
set (GEOGRAPHICLIB_LIB_TYPE_VAL 0)

set (PROJECT_STATIC_LIBRARIES geographic)
set (PROJECT_STATIC_DEFINITIONS -DGEOGRAPHICLIB_SHARED_LIB=0)

# (7) Set the default "real" precision.  This should probably be left
# at 2 (double).
set (GEOGRAPHICLIB_PRECISION 2 CACHE STRING
  "Precision: 1 = float, 2 = double, 3 = extended, 4 = quadruple, 5 = variable")
set_property (CACHE GEOGRAPHICLIB_PRECISION PROPERTY STRINGS 1 2 3 4 5)

set (LIBNAME geographic)

if (NOT MSVC)
  # Set the run time path for shared libraries for non-Windows machines.
  # (1) include link path for external packages (not needed with
  # GeographicLib because there are no external packages).  This only
  # makes sense for native builds.
  if (NOT CMAKE_CROSSCOMPILING)
    set (CMAKE_INSTALL_RPATH_USE_LINK_PATH TRUE)
  endif ()
  # (2) include installed path for GeographicLib.
  if (NOT APPLE)
    # Use relative path so that package is relocatable
    set (CMAKE_INSTALL_RPATH "\$ORIGIN/../lib${LIB_SUFFIX}")
  else ()
    if (CMAKE_VERSION VERSION_LESS 2.8.12)
      # Use absolute path with MacOSx
      set (CMAKE_INSTALL_NAME_DIR "${CMAKE_INSTALL_PREFIX}/lib${LIB_SUFFIX}")
    else ()
      # cmake 2.8.12 introduced a way to make the package relocatable.
      # See also INSTALL_RPATH property on the tools.
      set (CMAKE_MACOSX_RPATH ON)
    endif ()
  endif ()
endif ()

include (CheckTypeSize)
check_type_size ("long double" LONG_DOUBLE BUILTIN_TYPES_ONLY)
set (GEOGRAPHICLIB_HAVE_LONG_DOUBLE ${HAVE_LONG_DOUBLE})
include (TestBigEndian)
test_big_endian (GEOGRAPHICLIB_WORDS_BIGENDIAN)

# Make the compiler more picky.
if (MSVC)
  string (REGEX REPLACE "/W[0-4]" "" CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS}")
  # Turn on parallel builds for Visual Studio
  set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /W4 /MP")
else ()
  set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra")
  # check for C++11 support.  This flag is *not* propagated to clients
  # that use GeographicLib.  However, this is of no consequence.  When
  # the client code is being compiled (and the GeographicLib headers
  # being included), work-alike substitutions are used.
  include (CheckCXXCompilerFlag)
  # Not yet ready for C++14 (problem with MPFR C++)
  foreach (_F 11 0x)
    set (CXX11_FLAG "-std=c++${_F}")
    set (_T CXX11TEST${_F})
    check_cxx_compiler_flag (${CXX11_FLAG} ${_T})
    if (${_T})
      break ()
    endif ()
    unset (CXX11_FLAG)
  endforeach ()
  if (CXX11_FLAG)
    set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} ${CXX11_FLAG}")
  endif ()
endif ()

include (CheckCXXSourceCompiles)
set (CMAKE_REQUIRED_FLAGS "${CMAKE_CXX_FLAGS}")
# Check whether the C++11 math function: std::expm1, std::atanh,
# etc. are available.  This flag is *not* propagated to clients that use
# GeographicLib.  However, this is of no consequence.  When the client
# code is being compiled (and the GeographicLib headers being included),
# work-alike substitutions are used.
check_cxx_source_compiles (
  "#include <cmath>
int main() {
  int q;
  return int(std::hypot(3.0, 4.0) + std::expm1(0.5) +
             std::log1p(2.0) + std::asinh(10.0) +
             std::atanh(0.8) + std::cbrt(8.0) +
             std::fma(1.0, 2.0, 3.0) + std::remquo(100.0, 90.0, &q) +
             std::remainder(100.0, 90.0)) +
             std::isfinite(4.0) + std::isnan(0.0);
}\n" CXX11_MATH)
if (CXX11_MATH)
  add_definitions (-DGEOGRAPHICLIB_CXX11_MATH=1)
else ()
  add_definitions (-DGEOGRAPHICLIB_CXX11_MATH=0)
endif ()

# Check whether the C++11 static_assert macro is available.  This flag
# is *not* propagated to clients that use GeographicLib.  However, this
# is of no consequence.  When the client code is being compiled (and the
# GeographicLib headers being included), a work-alike substitution is
# used.
check_cxx_source_compiles (
  "#include <cmath>
int main() {
  static_assert(true, \"static assert test\");
  return 0;
}\n" CXX11_STATIC_ASSERT)
if (CXX11_STATIC_ASSERT)
  add_definitions (-DGEOGRAPHICLIB_HAS_STATIC_ASSERT=1)
else ()
  add_definitions (-DGEOGRAPHICLIB_HAS_STATIC_ASSERT=0)
endif ()

# Set the include directories.  Look in ${PROJECT_BINARY_DIR}/include
# first because that's where Config.h will be
include_directories ("${PROJECT_BINARY_DIR}/include" include)

if (GEOGRAPHICLIB_PRECISION EQUAL 1)
  message (WARNING "Compiling with floats which results in poor accuracy")
elseif (GEOGRAPHICLIB_PRECISION EQUAL 2)
  # This is the default
elseif (GEOGRAPHICLIB_PRECISION EQUAL 3)
  if (WIN32)
    message (WARNING
      "Cannot support long double on Windows, switching to double")
    set (GEOGRAPHICLIB_PRECISION 2)
  endif ()
elseif (GEOGRAPHICLIB_PRECISION EQUAL 4)
  if (CMAKE_CXX_COMPILER_ID STREQUAL GNU AND CXX11_MATH)
    # quad precision numbers appeared in Boost 1.54.
    find_package (Boost 1.54)
    if (Boost_FOUND)
      include_directories ("${Boost_INCLUDE_DIRS}")
      set (QUAD_LIBRARIES quadmath)
      if (NOT CMAKE_CXX_COMPILER_VERSION VERSION_LESS 4.8)
        # Enable Q suffix for quad precision in g++ 4.8 and later
        set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fext-numeric-literals")
      endif ()
      # Suppress the warnings that boost can spew out, e.g.,
      #   "typedef 'boost_concept_check905' locally defined but not used"
      set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wno-unused-local-typedefs")
    endif ()
  endif ()
  if (NOT QUAD_LIBRARIES)
    message (WARNING "Cannot support quad precision, switching to double")
    set (GEOGRAPHICLIB_PRECISION 2)
  endif ()
elseif (GEOGRAPHICLIB_PRECISION EQUAL 5)
  if (CXX11_MATH)
    # Install MPFR C++ version 3.6.2 or later from
    # http://www.holoborodko.com/pavel/mpfr and install mpreal.h in the
    # include directory.  NOTE: MPFR C++ is covered by the GPL; be sure
    # to abide by the terms of this license.
    #
    # For Linux, use system versions of mpfr and gmp.  For Apple, use
    # brew install mpfr.  Recent versions of mpfr (3.0 or later) work
    # fine.  For Windows, download MPFR-MPIR-x86-x64-MSVC2010.zip from
    # the MPFR C++ site and unpack in the top-level directory.  The
    # Windows build only works with GEOGRAPHICLIB_LIB_TYPE=STATIC.
    # NOTE: mpfr, gmp, and mpir are covered by the LGPL; be sure to
    # abide by the terms of this license.
    #
    # Need Visual Studio 12 2013 or later, g++ 4.5 or later; not sure
    # about clang.
    if (WIN32)
      if (MSVC AND NOT MSVC_VERSION LESS 1800)
        if (CMAKE_SIZEOF_VOID_P EQUAL 8)
          set (_ARCH x64)
        else ()
          set (_ARCH Win32)
        endif ()
        include_directories (mpfr_mpir_x86_x64_msvc2010/mpfr
          mpfr_mpir_x86_x64_msvc2010/mpir/dll/${_ARCH}/Release)
        # These are C libraries so it's OK to use release versions for
        # debug builds.  Also these work for later versions of Visual
        # Studio (specifically version 12).
        link_directories (mpfr_mpir_x86_x64_msvc2010/mpfr/dll/${_ARCH}/Release
          mpfr_mpir_x86_x64_msvc2010/mpir/dll/${_ARCH}/Release)
        set (MPFR_LIBRARIES mpfr mpir)
        # Suppress the myriad of "conditional expression is constant"
        # warnings
        set (CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /wd4127")
      endif ()
    else ()
      if (APPLE)
        include_directories (/usr/local/include)
        link_directories (/usr/local/lib)
      endif ()
      # g++ before 4.5 doesn't work (no explicit cast operator)
      if (NOT (CMAKE_CXX_COMPILER_ID STREQUAL GNU AND
            CMAKE_CXX_COMPILER_VERSION VERSION_LESS 4.5))
        set (MPFR_LIBRARIES mpfr gmp)
      endif ()
    endif ()
  endif ()
  if (NOT MPFR_LIBRARIES)
    message (WARNING "Cannot support mpfr, switching to double")
    set (GEOGRAPHICLIB_PRECISION 2)
  endif ()
endif ()

if (APPLE AND GEOGRAPHICLIB_PRECISION EQUAL 2)
  if (CMAKE_SYSTEM_PROCESSOR MATCHES "i.86" OR
      CMAKE_SYSTEM_PROCESSOR MATCHES "amd64" OR
      CMAKE_SYSTEM_PROCESSOR MATCHES "x86")
    set (CMAKE_OSX_ARCHITECTURES "x86_64")
  endif ()
endif ()

# Create a Config.h to expose system information to the compiler
configure_file (
  include/GeographicLib/Config.h.in
  include/GeographicLib/Config.h
  @ONLY )

if (WIN32)
  # static libaries
  set (CMAKE_ARCHIVE_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}/lib")
  # shared libraries
  set (CMAKE_LIBRARY_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}/lib")
  # executables and dlls
  set (CMAKE_RUNTIME_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}/bin")
endif ()

set_property (GLOBAL PROPERTY USE_FOLDERS ON)

# The list of subdirectories to process
add_subdirectory (src)
add_subdirectory (include/GeographicLib)
