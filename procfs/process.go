package procfs

type Process struct {
	PID      int
	Comm     string
	State    byte
	Priority int
	Nice     int

	PrivateKB uint64
	VSZBytes  uint64
	RSSBytes  uint64

	UTimeTicks uint64
	STimeTicks uint64
}
