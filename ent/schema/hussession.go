package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// HusSession holds the schema definition for the HusSession entity.
type HusSession struct {
	ent.Schema
}

// Fields of the HusSession.
func (HusSession) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).StructTag(`json:"sid,omitempty"`).Default(uuid.New).Unique(),
		field.UUID("uid", uuid.UUID{}),
		field.Bool("hld").Default(false),                              // holded
		field.Time("exp").Default(time.Now().Add(time.Hour * 24 * 7)), // expires at
		field.Time("iat").Default(time.Now),                           // issued at
	}
}

// Edges of the HusSession.
func (HusSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("hus_sessions").Unique().Field("uid").Required(),
	}
}
