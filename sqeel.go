package sqeel

import (
	"fmt"
	"reflect"
	"strings"
)

// Key stores information about a given SQL key and its creation attributes.
type Key struct {
	Name     string // Name of the key
	GoType   string // The Golang type of the key in string form
	SQLType  string // The SQL type of the key in string form
	SQLAttrs string // Additional SQL attributes (optional)
	Tag      string // Raw tag value
}

// SQLDefinition returns a keys SQL definition.
func (k *Key) SQLDefinition() string {
	sd := fmt.Sprintf("%s %s", k.Name, k.SQLType)
	if len(k.SQLAttrs) > 0 {
		sd += " " + k.SQLAttrs
	}
	return sd
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

		td.Keys = append(td.Keys, &Key{
			Name:     t.Name,
			GoType:   f.Type().String(),
			SQLType:  st.Type,
			SQLAttrs: st.Attrs,
			Tag:      tag,
		})
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

// CreationSchema returns this tables creation schema.
func (td *TableDescription) CreationSchema() string {
	q := fmt.Sprintf("CREATE TABLE %s (\n", td.Name)
	lines := []string{}
	for _, k := range td.Keys {
		lines = append(lines, k.SQLDefinition()+",")
	}
	lines = append(lines, fmt.Sprintf("PRIMARY KEY (%s)", td.PrimaryKey().Name))
	return q + strings.Join(lines, "\n") + ");"
}

// DeletionSchema returns this tables deletion schema.
func (td *TableDescription) DeletionSchema() string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;", td.Name)
}

// ValidateSchema returns this tables validation schema.
func (td *TableDescription) ValidateSchema() string {
	return fmt.Sprintf("SHOW TABLES LIKE '%s';", td.Name)
}
