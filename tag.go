package sqeel

import (
	"fmt"
	"strings"
)

// sqltag represents a string tag which will be interpreted into the SQL type,
// the attributes and if this is to be the tables primary key.  The only way
// to specify the primary key attribute is to set the second comma-separated
// value in the tag to an enable or "PRIMARYKEY".
type sqltag struct {
	NameOverride string
	Type         string
	Attrs        string
	IsPrimary    bool
	IsUnique     bool
}

// newtag returns a `*sqltag` from a string.
func newtag(tag string) sqltag {
	var stag sqltag
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		kv := strings.Split(part, ":")
		key, value := kv[0], ""
		if len(kv) > 1 {
			value = strings.Join(kv[1:], ":")
		}
		switch key {
		case "type":
			stag.Type = value
		case "is_primary", "primary", "primary_key", "primarykey":
			stag.IsPrimary = true
		case "attrs":
			stag.Attrs = value
		case "unique":
			stag.IsUnique = true
		case "name", "column_name":
			stag.NameOverride = value
		default:
			panic(fmt.Sprintf("invalid key found for sqeel tag (%s)", key))
		}
	}
	return stag
}
