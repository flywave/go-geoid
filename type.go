package geoid

type VerticalDatum uint32

const (
	HAE     VerticalDatum = 0
	EGM84   VerticalDatum = 1
	EGM96   VerticalDatum = 2
	EGM2008 VerticalDatum = 3
	UNKNOWN VerticalDatum = 4
)

func (v *VerticalDatum) ToString() string {
	switch *v {
	case HAE:
		return "HAE"
	case EGM84:
		return "EGM84"
	case EGM96:
		return "EGM96"
	case EGM2008:
		return "EGM2008"
	case UNKNOWN:
		return "UNKNOWN"
	}
	return "HAE"
}

type ConvertFlag int32

const (
	ELLIPSOIDTOGEOID ConvertFlag = -1
	NONE             ConvertFlag = 0
	GEOIDTOELLIPSOID ConvertFlag = 1
)
