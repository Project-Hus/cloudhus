package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ConnectedSessions holds the schema definition for the ConnectedSessions entity.
type ConnectedSessions struct {
	ent.Schema
}

func (ConnectedSessions) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Options: "ENGINE=MEMORY",
		},
	}
}

// Fields of the ConnectedSessions.
func (ConnectedSessions) Fields() []ent.Field {
	return []ent.Field{
		// hus session id
		field.UUID("hsid", uuid.UUID{}),
		// subservice's auth server origin
		field.String("service_origin"),
		// connected session id from subservice
		field.UUID("csid", uuid.UUID{}),
	}
}

// Edges of the ConnectedSessions.
func (ConnectedSessions) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("hus_session", HusSession.Type).Unique().Ref("connected_sessions").Field("hsid").Required(),
	}
}
