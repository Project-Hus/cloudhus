// Code generated by ent, DO NOT EDIT.

package hussession

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the hussession type in the database.
	Label = "hus_session"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTid holds the string denoting the tid field in the database.
	FieldTid = "tid"
	// FieldIat holds the string denoting the iat field in the database.
	FieldIat = "iat"
	// FieldExp holds the string denoting the exp field in the database.
	FieldExp = "exp"
	// FieldPreserved holds the string denoting the preserved field in the database.
	FieldPreserved = "preserved"
	// FieldUID holds the string denoting the uid field in the database.
	FieldUID = "uid"
	// FieldSignedAt holds the string denoting the signed_at field in the database.
	FieldSignedAt = "signed_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeConnectedSession holds the string denoting the connected_session edge name in mutations.
	EdgeConnectedSession = "connected_session"
	// Table holds the table name of the hussession in the database.
	Table = "hus_sessions"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "hus_sessions"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "uid"
	// ConnectedSessionTable is the table that holds the connected_session relation/edge.
	ConnectedSessionTable = "connected_sessions"
	// ConnectedSessionInverseTable is the table name for the ConnectedSession entity.
	// It exists in this package in order to avoid circular dependency with the "connectedsession" package.
	ConnectedSessionInverseTable = "connected_sessions"
	// ConnectedSessionColumn is the table column denoting the connected_session relation/edge.
	ConnectedSessionColumn = "hsid"
)

// Columns holds all SQL columns for hussession fields.
var Columns = []string{
	FieldID,
	FieldTid,
	FieldIat,
	FieldExp,
	FieldPreserved,
	FieldUID,
	FieldSignedAt,
	FieldUpdatedAt,
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

var (
	// DefaultTid holds the default value on creation for the "tid" field.
	DefaultTid func() uuid.UUID
	// DefaultIat holds the default value on creation for the "iat" field.
	DefaultIat func() time.Time
	// DefaultExp holds the default value on creation for the "exp" field.
	DefaultExp int64
	// UpdateDefaultExp holds the default value on update for the "exp" field.
	UpdateDefaultExp func() int64
	// DefaultPreserved holds the default value on creation for the "preserved" field.
	DefaultPreserved bool
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
