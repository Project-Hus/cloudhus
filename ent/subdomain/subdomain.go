// Code generated by ent, DO NOT EDIT.

package subdomain

const (
	// Label holds the string label denoting the subdomain type in the database.
	Label = "subdomain"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldServiceID holds the string denoting the service_id field in the database.
	FieldServiceID = "service_id"
	// FieldSubdomain holds the string denoting the subdomain field in the database.
	FieldSubdomain = "subdomain"
	// FieldRole holds the string denoting the role field in the database.
	FieldRole = "role"
	// EdgeService holds the string denoting the service edge name in mutations.
	EdgeService = "service"
	// Table holds the table name of the subdomain in the database.
	Table = "subdomains"
	// ServiceTable is the table that holds the service relation/edge.
	ServiceTable = "subdomains"
	// ServiceInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceInverseTable = "services"
	// ServiceColumn is the table column denoting the service relation/edge.
	ServiceColumn = "service_id"
)

// Columns holds all SQL columns for subdomain fields.
var Columns = []string{
	FieldID,
	FieldServiceID,
	FieldSubdomain,
	FieldRole,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}
