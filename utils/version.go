package utils

import (
	"fmt"
)

type VersionInt struct {
	Major    int
	Minor    int
	Revision int
	Build    string
}

func NewVersionInt(verstr string, sepStr string) VersionInt {
	var fmtstr string
	if sepStr == "" {
		fmtstr = "%d.%d.%d"
	} else {
		fmtstr = "%d.%d.%d" + sepStr + "%s"
	}
	var v VersionInt
	fmt.Sscanf(verstr, fmtstr, &v.Major, &v.Minor, &v.Revision, &v.Build)
	return v
}

func (v *VersionInt) GreatThan(nv *VersionInt) bool {
	if v.Major > nv.Major {
		return true
	} else if v.Major < nv.Major {
		return false
	}

	if v.Minor > nv.Minor {
		return true
	} else if v.Minor < nv.Minor {
		return false
	}

	return v.Revision > nv.Revision
}

func (v *VersionInt) Equal(nv *VersionInt) bool {
	return *v == *nv
}

func (v *VersionInt) Compatible(nv *VersionInt) bool {
	return v.Major == nv.Major && v.Minor == nv.Minor
}

func (v *VersionInt) LessThan(nv *VersionInt) bool {
	if v.Major < nv.Major {
		return true
	} else if v.Major > nv.Major {
		return false
	}

	if v.Minor < nv.Minor {
		return true
	} else if v.Minor > nv.Minor {
		return false
	}

	return v.Revision < nv.Revision
}
