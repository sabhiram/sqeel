package sqeel

import "strings"

// sqltag represents a string tag which will be interpreted into the SQL type,
// the attributes and if this is to be the tables primary key.  The only way
// to specify the primary key attribute is to set the second comma-separated
// value in the tag to an enable or "PRIMARYKEY".
type sqltag struct {
	Type      string
	Attrs     string
	IsPrimary bool
}

// newtag returns a `*sqltag` from a string.
func newtag(tag string) sqltag {
	parts := strings.Split(tag, ",")
	t, attrs := parts[0], ""
	tparts := strings.Split(t, " ")
	if len(tparts) >= 1 {
		t = tparts[0]
		attrs = strings.Join(tparts[1:], " ")
	}

	pk := false
	if len(parts) >= 2 {
		pkv := strings.ToLower(parts[1])
		if pkv == "1" || pkv == "true" || pkv == "primary" || pkv == "primarykey" {
			pk = true
		}
	}

	return sqltag{
		Type:      t,
		Attrs:     attrs,
		IsPrimary: pk,
	}
}
