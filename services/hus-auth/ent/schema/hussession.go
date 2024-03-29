package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// HusSession holds the schema definition for the HusSession entity.
type HusSession struct {
	ent.Schema
}

func (HusSession) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Options: "ENGINE=MEMORY",
		},
	}
}

func GetHstExp() int64 {
	return time.Now().Add(time.Hour * 48).Unix() // 10 min basically
}

// Fields of the HusSession.
func (HusSession) Fields() []ent.Field {
	return []ent.Field{
		// ID of the Hus session
		field.UUID("id", uuid.UUID{}).StructTag(`json:"sid,omitempty"`).Default(uuid.New).Unique(),
		// Session's temporary ID for rotation
		field.UUID("tid", uuid.UUID{}).Default(uuid.New),
		// issued at
		field.Time("iat").Default(time.Now),
		// expired at
		field.Int64("exp").Default(GetHstExp()).UpdateDefault(GetHstExp),
		// in preserved mode, the session token is extended by a week each time the user is redirected to cloudhus.
		// but tid is rotated each time.
		field.Bool("preserved").Default(false), // preserved
		// User ID for the case the user signed in.
		field.Uint64("uid").Optional().Nillable(),
		field.Time("signed_at").Optional().Nillable(),

		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the HusSession.
func (HusSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("hus_sessions").Unique().Field("uid"),

		edge.To("connected_session", ConnectedSession.Type),
	}
}
