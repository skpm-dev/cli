package version

import (
	"fmt"
	"strconv"
	"strings"
)

type BumpType int

const (
	BumpPatch BumpType = iota
	BumpMinor
	BumpMajor
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func Parse(v string) (*Version, error) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid version format %q, expected major.minor.patch", v)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %w", err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %w", err)
	}

	return &Version{Major: major, Minor: minor, Patch: patch}, nil
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) Bump(t BumpType) *Version {
	switch t {
	case BumpMajor:
		return &Version{Major: v.Major + 1, Minor: 0, Patch: 0}
	case BumpMinor:
		return &Version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}
	case BumpPatch:
		return &Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
	default:
		return v
	}
}
