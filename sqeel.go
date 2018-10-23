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
	Name       string            // Name of the table
	PrimaryKey string            // PrimaryKey variable name
	Keys       []*Key            // List of keys (ordered)
	Attrs      map[string]string // Custom attrs that the user wants to access

	imports []string // local list of imports we have encountered
}

// DescribeTable accepts a table name, an interface that employs the sqeel tags
// as well as any foreign key setup.  The attributes are a map of string ->
// string which is opaque to this implementation.  Returns a TableDescription
// instance.  Panics on any errors.
func DescribeTable(name string, v interface{}, fks, attrs map[string]string) *TableDescription {
	td := &TableDescription{
		Name:       name,
		PrimaryKey: "",
		Keys:       []*Key{},
		Attrs:      attrs,

		imports: []string{},
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
		key := &Key{
			Name:     t.Name,
			GoType:   f.Type().String(),
			SQLType:  st.Type,
			SQLAttrs: st.Attrs,
			Tag:      tag,
		}
		tn := key.GoType
		switch {
		case tn == "time.Time":
			td.appendImport("time")
		case strings.Contains(tn, "."):
			fmt.Printf("Unknown type %s\n", tn)
			panic("bad type in table import")
		}
		td.Keys = append(td.Keys, key)
	}

	if td.PrimaryKey == "" {
		panic("Unable to find primary key field")
	}
	return td
}

func (td *TableDescription) appendImport(imp string) {
	for _, v := range td.imports {
		if v == imp {
			return
		}
	}
	td.imports = append(td.imports, imp)
}

// Imports returns a unique list of all the import types that we have included
// so far.
func (td *TableDescription) Imports() []string {
	if td == nil {
		return nil
	}
	return td.imports
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
