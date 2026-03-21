package zerorat

import (
	"bytes"
	"encoding/json"
	"errors"
)

var (
	// errNilRatJSONTarget reports attempts to unmarshal JSON into a nil Rat receiver.
	errNilRatJSONTarget = errors.New("cannot unmarshal JSON into nil *Rat")
	// Production assertions keep the JSON contract tied to the standard interfaces.
	_ json.Marshaler   = Rat{}
	_ json.Unmarshaler = (*Rat)(nil)
)

// MarshalJSON encodes Rat as a normalized JSON number.
// It rejects invalid and non-terminating values rather than changing them silently.
func (r Rat) MarshalJSON() ([]byte, error) {
	decimal, err := r.ToDecimalString()
	if err != nil {
		return nil, err
	}

	return []byte(decimal), nil
}

// UnmarshalJSON decodes Rat from a JSON number, quoted decimal text, or null.
// It invalidates the receiver for null and for failed decodes so stale values do not survive bad payloads.
func (r *Rat) UnmarshalJSON(data []byte) error {
	if r == nil {
		return errNilRatJSONTarget
	}

	raw := bytes.TrimSpace(data)
	if string(raw) == "null" {
		r.Invalidate()
		return nil
	}

	text := string(raw)
	if len(raw) > 0 && raw[0] == '"' {
		var quoted string
		if err := json.Unmarshal(raw, &quoted); err != nil {
			r.Invalidate()
			return err
		}
		text = quoted
	}

	parsed, err := NewFromDecimalString(text)
	if err != nil {
		r.Invalidate()
		return err
	}

	*r = parsed
	return nil
}
