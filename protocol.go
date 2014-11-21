package keep

import (
	"bytes"
	"fmt"
	"time"
)

var (
	nodekind = []byte("notes#node")
	tskind   = []byte("notes#timestamps")

	notetype = []byte("NOTE")
	itemtype = []byte("LIST_ITEM")

	tszero = []byte("1970-01-01T00:00:00.000Z")
)

const (
	tsLayout = "2006-01-02T15:04:05.000Z"
)

// note represents a textual entry in google keep.
type note struct {
	node
	Type noteType `json:"type"`
}

// node represents an identity of an item of data in google keep.
type node struct {
	ID         string     `json:"id"`
	Kind       nodeKind   `json:"kind"`
	ParentId   string     `json:"parentId"`
	SortValue  int        `json:"sortValue"`
	Timestamps timestamps `json:"timestamps"`
}

// timestamps is a collection of time-related data about a node.
type timestamps struct {
	Kind       tsKind    `json:"kind"`
	Created    timestamp `json:"created"`
	Deleted    timestamp `json:"deleted"`
	Trashed    timestamp `json:"trashed"`
	Updated    timestamp `json:"updated"`
	UserEdited timestamp `json:"userEdited"`
}

// timestamp is a time that serializes to a string where time zero = 1970.
type timestamp time.Time

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (t timestamp) MarshalJSON() ([]byte, error) {
	tm := time.Time(t)
	if t.IsZero() {
		return tszero, nil
	}
	return []byte(tm.Format(tsLayout))
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (t *timestamp) UnmarshalJSON(b []byte) error {
	if bytes.Equal(tszero, b) {
		t = timestamp{}
		return nil
	}
	tm, err := time.Parse(tsLayout, string(b))
	if err != nil {
		return err
	}
	*t = timestamp(tm)
}

// tsKind outputs the correct json for the kind field of the timestamps struct.
type tsKind struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (tsKind) MarshalJSON() ([]byte, error) {
	return tskind
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*tsKind) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(tskind, b) {
		return fmt.Errorf("expected %q got %q", tskind, b)
	}
	return nil
}

// nodeKind outputs the correct json for the kind field of the timestamps struct.
type nodeKind struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (nodeKind) MarshalJSON() ([]byte, error) {
	return nodekind
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*nodeKind) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(nodekind, b) {
		return fmt.Errorf("expected %q got %q", nodekind, b)
	}
	return nil
}

// noteType outputs the correct json for the type field of the note struct.
type noteType struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (noteType) MarshalJSON() ([]byte, error) {
	return notetype
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*noteType) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(notetype, b) {
		return fmt.Errorf("expected %q got %q", notetype, b)
	}
	return nil
}
