package eventStore

type ConcurrencyMode int

const (
	ANY ConcurrencyMode = iota
	NEW_STREAM
	EXISTING_STREAM
	EXPECTING_SEQ_NUM
)

type EventEnvelope struct {
	SeqNum    int64  `json:"n"`
	Timestamp int64  `json:"ts"`
	EventType string `json:"et"`
	Data      []byte `json:"d"`
}

type EventStore interface {
	GetNamespaces() ([]string, error)

	GetStreams(ns string) ([]string, error)

	NamespaceExists(ns string) (bool, error)

	StreamExists(ns string, streamId string) (bool, error)

	WriteEvent(ns string, streamId string,
		cMode ConcurrencyMode, expected int64, e *EventEnvelope) (int64, error)

	WriteBatch(ns string, streamId string,
		cMode ConcurrencyMode, expected int64, es []EventEnvelope) (int64, error)

	GetEvent(ns string, streamId string,
		seqNum int64) (*EventEnvelope, error)

	GetEventRange(ns string, streamId string,
		starting int64, ending int64) ([]EventEnvelope, error)
}
