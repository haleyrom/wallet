package servant

// Configure Configure
type Configure interface {
	Init(p string) error
}
