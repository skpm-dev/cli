package version

import "github.com/Masterminds/semver/v3"

type BumpType int

const (
	BumpPatch BumpType = iota
	BumpMinor
	BumpMajor
)

func Parse(v string) (*semver.Version, error) {
	return semver.NewVersion(v)
}

func Bump(v *semver.Version, t BumpType) semver.Version {
	switch t {
	case BumpMajor:
		return v.IncMajor()
	case BumpMinor:
		return v.IncMinor()
	default:
		return v.IncPatch()
	}
}
