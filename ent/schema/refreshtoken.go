package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// RefreshToken holds the schema definition for the RefreshToken entity.
type RefreshToken struct {
	ent.Schema
}

// Fields of the RefreshToken.
func (RefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("uid"), // User ID
		field.Bool("revoked").Default(false),
		field.Time("last_used_at").Default(time.Time{}), //  January 1, year 1, 00:00:00.000000000 UTC
		field.Time("created_at").Default(time.Now()),
	}
}

// Edges of the RefreshToken.
func (RefreshToken) Edges() []ent.Edge {
	return nil
}
