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
	Name       string // Name of the table
	PrimaryKey string // PrimaryKey variable name
	Keys       []*Key // List of keys (ordered)
}

// DescribeTable accepts a table name, an interface that employs the sqeel tags
// as well as any foreign key setup.  The attributes are a map of string ->
// string which is opaque to this implementation.  Returns a TableDescription
// instance.  Panics on any errors.
func DescribeTable(name string, v interface{}, fks map[string]string) *TableDescription {
	td := &TableDescription{
		Name:       name,
		PrimaryKey: "",
		Keys:       []*Key{},
	}

	e := reflect.ValueOf(v)
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		t := e.Type().Field(i)
		tag := t.Tag.Get("sqeel")

		st := newtag(tag)
		if st.IsPrimary {
			td.PrimaryKey = t.Name
		}

		td.Keys = append(td.Keys, &Key{
			Name:     t.Name,
			GoType:   f.Type().String(),
			SQLType:  st.Type,
			SQLAttrs: st.Attrs,
			Tag:      tag,
		})
	}

	if len(td.PrimaryKey) <= 0 {
		panic("Unable to find primary key field")
	}
	return td
}

// CreationSchema returns this tables creation schema.
func (td *TableDescription) CreationSchema() string {
	q := fmt.Sprintf("CREATE TABLE %s (\n", td.Name)
	lines := []string{}
	for _, k := range td.Keys {
		lines = append(lines, k.SQLDefinition()+",")
	}
	lines = append(lines, fmt.Sprintf("PRIMARY KEY (%s)", td.PrimaryKey))
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
