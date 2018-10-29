package sqeel

import (
	"fmt"
	"reflect"
	"strings"
)

// Key stores information about a given SQL key and its creation attributes.
type Key struct {
	Name       string // Name of the key
	ColumnName string // Optional, override with `name` tag
	GoType     string // The Golang type of the key in string form
	SQLType    string // The SQL type of the key in string form
	SQLAttrs   string // Additional SQL attributes (optional)
	Tag        string // Raw tag value
	IsPrimary  bool   // If this is the primary key of the table
	IsUnique   bool   // If this is a unique key in the table
}

// SQLDefinition returns a key's SQL definition.
func (k *Key) SQLDefinition() string {
	sd := fmt.Sprintf("%s %s", k.Name, k.SQLType)
	if len(k.SQLAttrs) > 0 {
		sd += " " + k.SQLAttrs
	}
	return sd
}

// GoName returns the name of the key as created in the source struct.  This is
// the go type.
func (k *Key) GoName() string {
	return k.Name
}

// SQLName returns the name of the key as created in the source table.  This is
// the sql column name.
func (k *Key) SQLName() string {
	if len(k.ColumnName) > 0 {
		return k.ColumnName
	}
	return ToSnakeCase(k.Name)
}

// TableDescription stores a bunch of interesting things about a SQL table.
type TableDescription struct {
	Name          string // Name of the table
	Keys          []*Key // List of keys (ordered)
	PrimaryKeyIdx int    // Index in the `Keys` list of the primary key
}

// DescribeTable accepts a table name, an interface that employs the sqeel tags
// as well as any foreign key setup.  The attributes are a map of string ->
// string which is opaque to this implementation.  Returns a TableDescription
// instance.  Panics on any errors.
func DescribeTable(name string, v interface{}, fks map[string]string) *TableDescription {
	td := &TableDescription{
		Name:          name,
		Keys:          []*Key{},
		PrimaryKeyIdx: -1,
	}

	e := reflect.ValueOf(v)
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		t := e.Type().Field(i)
		tag := t.Tag.Get("sqeel")

		st := newtag(tag)
		if st.IsPrimary {
			td.PrimaryKeyIdx = i
		}

		k := &Key{
			Name:      t.Name,
			GoType:    f.Type().String(),
			SQLType:   st.Type,
			SQLAttrs:  st.Attrs,
			Tag:       tag,
			IsPrimary: st.IsPrimary,
			IsUnique:  st.IsUnique,
		}

		if len(st.NameOverride) > 0 {
			k.ColumnName = st.NameOverride
		}
		td.Keys = append(td.Keys, k)
	}

	if td.PrimaryKeyIdx < 0 {
		panic("primary key not specified in sqeel structure")
	}
	return td
}

// PrimaryKey returns the key that is tagged as the tables primary key.
func (td *TableDescription) PrimaryKey() *Key {
	return td.Keys[td.PrimaryKeyIdx]
}

// KeyNames returns a list of the key names.
func (td *TableDescription) KeyNames() []string {
	kns := []string{}
	for _, k := range td.Keys {
		kns = append(kns, k.Name)
	}
	return kns
}

// CreateTableQuery returns this tables creation query.
func (td *TableDescription) CreateTableQuery() string {
	q := fmt.Sprintf("CREATE TABLE %s (\n", td.Name)
	lines := []string{}
	keyls := []string{fmt.Sprintf("PRIMARY KEY (%s)", td.PrimaryKey().Name)}
	for _, k := range td.Keys {
		lines = append(lines, k.SQLDefinition())
		if k.IsUnique {
			keyls = append(keyls, fmt.Sprintf("UNIQUE KEY (%s)", k.Name))
		}
	}

	lines = append(lines, keyls...)
	return q + strings.Join(lines, ",\n") + ");"
}

// DeleteTableQuery returns this tables deletion query.
func (td *TableDescription) DeleteTableQuery() string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;", td.Name)
}

// ValidateTableQuery returns this tables validation query.
func (td *TableDescription) ValidateTableQuery() string {
	return fmt.Sprintf("SHOW TABLES LIKE '%s';", td.Name)
}
