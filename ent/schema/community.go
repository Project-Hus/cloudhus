package schema

import (
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Group holds the schema definition for the Group entity.
type Community struct {
	ent.Schema
}

// Fields of the Group.
func (Community) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Default(uuid.New().String()),
		field.String("name").
			Match(regexp.MustCompile("[a-zA-Z_]+$")),
	}
}

// Edges of the Group.
func (Community) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
