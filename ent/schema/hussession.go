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
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.Time("expired_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the HusSession.
func (HusSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("hus_sessions").Unique(),
	}
}
