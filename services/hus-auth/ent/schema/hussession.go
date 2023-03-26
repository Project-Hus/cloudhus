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
		field.UUID("id", uuid.UUID{}).StructTag(`json:"sid,omitempty"`).Default(uuid.New).Unique(), // sid
		field.UUID("tid", uuid.UUID{}).Default(uuid.New).Unique(),                                  // tid
		field.Time("iat").Default(time.Now),                                                        // issued at
		// if exp is nil, the session expires when the brwoser's session ends.
		field.Bool("preserved").Default(false), // preserved
		field.UUID("uid", uuid.UUID{}),         // uear id
	}
}

// Edges of the HusSession.
func (HusSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("hus_sessions").Unique().Field("uid").Required(),
	}
}
