package registry

type Version struct {
	Major int
	Minor int
	Patch int
}

func NewVersion(major, minor, patch int) Version {
	return Version{major, minor, patch}
}
