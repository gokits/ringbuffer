package ringstore

type Cursor struct {
	next    uint64
	initial bool

	tmpNext    uint64
	tmpInitial bool
}

func NewCursor() *Cursor {
	return &Cursor{}
}

func (c *Cursor) Commit() {
	//TODO check if commit in dirty state
	c.next = c.tmpNext
	c.initial = c.tmpInitial
}

func (c *Cursor) Rollback() {
	//TODO check if commit in dirty state
	c.tmpInitial = c.initial
	c.tmpNext = c.next
}

type RingStoreGetState int

func (rs RingStoreGetState) String() string {
	switch rs {
	case Behind:
		return "Behind"
	case Match:
		return "Match"
	case Beyond:
		return "Beyond"
	case Empty:
		return "Empty"
	default:
		return "Not Defined"
	}
}

const (
	Behind RingStoreGetState = iota
	Match
	Beyond
	Empty
)

type RingBuffer interface {
	Get(Cursor) (interface{}, RingStoreGetState, error)
	Append(interface{}) error
	Close() error
}
