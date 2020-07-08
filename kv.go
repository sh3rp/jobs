package jobs

type K string
type V struct {
	Type  string
	Value interface{}
}

type KV struct {
	K K
	V V
}
