package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ConnectedSession holds the schema definition for the ConnectedSessions entity.
type ConnectedSession struct {
	ent.Schema
}

func (ConnectedSession) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Options: "ENGINE=MEMORY",
		},
	}
}

// Fields of the ConnectedSession.
func (ConnectedSession) Fields() []ent.Field {
	return []ent.Field{
		// hus session id
		field.UUID("hsid", uuid.UUID{}),
		// subservice name
		field.String("service"),
		// connected session id from subservice
		field.UUID("csid", uuid.UUID{}),
	}
}

// Edges of the ConnectedSessions.
func (ConnectedSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("hus_session", HusSession.Type).Unique().Ref("connected_session").Field("hsid").Required(),
	}
}
