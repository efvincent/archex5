package esErrors

import (
	"fmt"
)

func NewStreamExists(sid string) *ESError {
	return &ESError{STREAM_EXISTS, sid, 0, 0}
}

func NewStreamDoesNotExist(sid string) *ESError {
	return &ESError{STREAM_DOES_NOT_EXIST, sid, 0, 0}
}

func NewSeqExpectedErr(sid string, expected int64, actual int64) *ESError {
	return &ESError{SEQ_NUM_EXPECTATION_FAILED, sid, expected, actual}
}

type ESErrorCode int

const (
	STREAM_EXISTS ESErrorCode = iota
	STREAM_DOES_NOT_EXIST
	SEQ_NUM_EXPECTATION_FAILED
)

type ESError struct {
	ErrCode  ESErrorCode
	StreamId string
	Expected int64
	Actual   int64
}

func (e ESError) Error() string {
	switch e.ErrCode {
	case STREAM_EXISTS:
		return fmt.Sprintf("Stream (%s) already exists", e.StreamId)
	case STREAM_DOES_NOT_EXIST:
		return fmt.Sprintf("Stream (%s) does not exist", e.StreamId)
	case SEQ_NUM_EXPECTATION_FAILED:
		if e.Actual == -1 {
			return fmt.Sprintf("Stream (%s) Expected last sequnce = %v, actual stream contains no events",
				e.StreamId, e.Expected)
		} else {
			return fmt.Sprintf("Stream (%s) Expected last sequnce = %v, actual last sequence = %v",
				e.StreamId, e.Expected, e.Actual)
		}
	default:
		return fmt.Sprintf("Unknown error code: %v", e.ErrCode)
	}
}

func (e ESError) String() string {
	return e.Error()
}
