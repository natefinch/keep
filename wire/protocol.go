package wire

import (
	"bytes"
	"fmt"
	"time"
)

var (
	// wire values for kind fields
	nodekind = []byte(`"notes#node"`)
	tskind   = []byte(`"notes#timestamps"`)

	// wire values for nodes' (static) type field
	notetype = []byte(`"NOTE"`)
	listtype = []byte(`"LIST"`)
	itemtype = []byte(`"LIST_ITEM"`)

	// wire value for the zero time
	tszero = []byte(`"1970-01-01T00:00:00.000Z"`)

	// wire values for reminder.state
	notDismissedVal = []byte(`"INITIAL"`)
	dismissedVal    = []byte(`"DISMISSED"`)

	// wire values for reminder.period
	morningVal   = []byte(`"MORNING"`)
	afternoonVal = []byte(`"AFTERNOON"`)
	eveningVal   = []byte(`"EVENING"`)
	nightVal     = []byte(`"NIGHT"`)

	// wire values for node colors
	defaultVal = []byte(`"DEFAULT"`)
	redVal     = []byte(`"RED"`)
	orangeVal  = []byte(`"ORANGE"`)
	yellowVal  = []byte(`"YELLOW"`)
	greenVal   = []byte(`"GREEN"`)
	tealVal    = []byte(`"TEAL"`)
	blueVal    = []byte(`"BLUE"`)
	grayVal    = []byte(`"GRAY"`)
)

const tsLayout = "2006-01-02T15:04:05.000Z"

// These are values for the Period field of Reminders.
const (
	SpecificTime Period = iota
	Morning             // 9am
	Afternoon           // 1pm
	Evening             // 5pm
	Night               // 8pm
)

// These are values for the Color field of notes and lists.
const (
	DefaultColor Color = iota
	Red
	Orange
	Yellow
	Green
	Teal
	Blue
	Gray
)

// Note represents a textual entry in google keep.  Notes have a single child
// node that contains the text of the note.
type Note struct {
	ParentNode
	Type NoteType `json:"type"`
}

// List represents a list of items in google keep.  Lists have one or more child
// nodes that are the items in the list.
type List struct {
	ParentNode
	Type ListType `json:"type"`
}

// Item represents a textual entry in google keep.
type Item struct {
	Node
	Type    ItemType `json:"type"`
	Checked bool     `json:"checked"`
	Text    string   `json:"text"`
}

// ParentNode is the base representation for the top-level node for a list or
// note.
type ParentNode struct {
	Node
	Title    string `json:"title"`
	Archived bool   `json:"isArchived"`
	Color
}

// Node represents an identity of an item of data in google keep.
type Node struct {
	ID         string     `json:"id"`
	Kind       NodeKind   `json:"kind"`
	ParentID   string     `json:"parentId"`
	SortValue  int        `json:"sortValue"`
	Timestamps Timestamps `json:"timestamps"`
}

// Timestamps is a collection of time-related data about a node.
type Timestamps struct {
	Kind       TSKind    `json:"kind"`
	Created    Timestamp `json:"created"`
	Deleted    Timestamp `json:"deleted"`
	Trashed    Timestamp `json:"trashed"`
	Updated    Timestamp `json:"updated"`
	UserEdited Timestamp `json:"userEdited"`
}

// Reminder contains a time and message to tell the user about a note or list.
type Reminder struct {
	Dismissed   Dismissed `json:"state"`
	Description string    `json:"description"`
	Time
}

// Time represents when a Reminder should notify the user.  The year, month, and
// day are always specified.  The time of day is specified either by a generic
// Period (i.e. Morning, Afternoon, etc), or by a specific time stored in the
// Hour, Minute, and Second fields.
type Time struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`

	Period Period `json:"period,omitempty"`

	Hour   int `json:"hour,omitempty"`
	Minute int `json:"minute,omitempty"`
	Second int `json:"second,omitempty"`
}

// Time returns the time.Time value represented by this struct.
func (t *Time) Time() time.Time {
	var h, m, s int
	switch t.Period {
	case SpecificTime:
		h, m, s = t.Hour, t.Minute, t.Second
	case Morning:
		h = 9
	case Afternoon:
		h = 13
	case Evening:
		h = 17
	case Night:
		h = 20
	}
	return time.Date(t.Year, time.Month(t.Month), t.Day, h, m, s, 0, time.Now().Location())
}

// Dismissed is a boolean type that serializes to DISMISSED or INITIAL
type Dismissed bool

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (d Dismissed) MarshalJSON() ([]byte, error) {
	if d {
		return dismissedVal, nil
	}
	return notDismissedVal, nil
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (d *Dismissed) UnmarshalJSON(b []byte) error {
	if bytes.Equal(dismissedVal, b) {
		*d = true
		return nil
	}
	if bytes.Equal(notDismissedVal, b) {
		*d = true
		return nil
	}

	return fmt.Errorf("expected %q or %q, got %q", dismissedVal, notDismissedVal, b)
}

// Period defines a broad timespan for reminders.
type Period int

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (p Period) MarshalJSON() ([]byte, error) {
	// note, Period should always be omitempty, so we don't support serializing
	// the default value here.
	switch p {
	case Morning:
		return morningVal, nil
	case Afternoon:
		return afternoonVal, nil
	case Evening:
		return eveningVal, nil
	case Night:
		return nightVal, nil
	}
	return nil, fmt.Errorf("unsupported period value %d", p)
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (p *Period) UnmarshalJSON(b []byte) error {
	// note, Period should always be omitempty, so we don't support
	// deserializing the default value here.
	switch {
	case bytes.Equal(morningVal, b):
		*p = Morning
	case bytes.Equal(afternoonVal, b):
		*p = Afternoon
	case bytes.Equal(eveningVal, b):
		*p = Evening
	case bytes.Equal(nightVal, b):
		*p = Night
	default:
		return fmt.Errorf("unexpected Period value %q", b)
	}
	return nil
}

// Color represents the color of a list or node.
type Color int

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (c Color) MarshalJSON() ([]byte, error) {
	switch c {
	case DefaultColor:
		return defaultVal, nil
	case Red:
		return redVal, nil
	case Orange:
		return orangeVal, nil
	case Yellow:
		return yellowVal, nil
	case Green:
		return greenVal, nil
	case Teal:
		return tealVal, nil
	case Blue:
		return blueVal, nil
	case Gray:
		return grayVal, nil
	}
	return nil, fmt.Errorf("unsupported color value %d", c)
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (c *Color) UnmarshalJSON(b []byte) error {
	switch {
	case bytes.Equal(defaultVal, b):
		*c = DefaultColor
	case bytes.Equal(redVal, b):
		*c = Red
	case bytes.Equal(orangeVal, b):
		*c = Orange
	case bytes.Equal(yellowVal, b):
		*c = Yellow
	case bytes.Equal(greenVal, b):
		*c = Green
	case bytes.Equal(tealVal, b):
		*c = Teal
	case bytes.Equal(blueVal, b):
		*c = Blue
	case bytes.Equal(grayVal, b):
		*c = Gray
	default:
		return fmt.Errorf("unexpected Color value %q", b)
	}
	return nil
}

// Timestamp is a time that serializes to a string where time zero = 1970.
type Timestamp time.Time

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return tszero, nil
	}
	return []byte(tm.Format(tsLayout)), nil
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	if bytes.Equal(tszero, b) {
		*t = Timestamp{}
		return nil
	}
	tm, err := time.Parse(tsLayout, string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(tm)
	return nil
}

// TSKind outputs the correct json for the kind field of the timestamps struct.
type TSKind struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (TSKind) MarshalJSON() ([]byte, error) {
	return tskind, nil
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*TSKind) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(tskind, b) {
		return fmt.Errorf("expected %q got %q", tskind, b)
	}
	return nil
}

// NodeKind outputs the correct json for the kind field of the timestamps struct.
type NodeKind struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (NodeKind) MarshalJSON() ([]byte, error) {
	return nodekind, nil
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*NodeKind) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(nodekind, b) {
		return fmt.Errorf("expected %q got %q", nodekind, b)
	}
	return nil
}

// NoteType outputs the correct json for the type field of the note struct.
type NoteType struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (NoteType) MarshalJSON() ([]byte, error) {
	return notetype, nil
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*NoteType) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(notetype, b) {
		return fmt.Errorf("expected %q got %q", notetype, b)
	}
	return nil
}

// ListType outputs the correct json for the type field of the list struct.
type ListType struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (ListType) MarshalJSON() ([]byte, error) {
	return listtype, nil
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*ListType) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(listtype, b) {
		return fmt.Errorf("expected %q got %q", listtype, b)
	}
	return nil
}

// ItemType outputs the correct json for the type field of the list struct.
type ItemType struct{}

// MarshalJSON implements json.Marhsaler.MarshalJSON.
func (ItemType) MarshalJSON() ([]byte, error) {
	return itemtype, nil
}

// UnmarshalJSON implements json.Unmarhsaler.UnmarshalJSON.
func (*ItemType) UnmarshalJSON(b []byte) error {
	if !bytes.Equal(itemtype, b) {
		return fmt.Errorf("expected %q got %q", itemtype, b)
	}
	return nil
}
