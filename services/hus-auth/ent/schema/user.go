package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).StructTag(`json:"uid,omitempty"`).Default(uuid.New).Unique(),

		field.Enum("provider").Values("hus", "google"),

		field.String("google_sub").Unique(),

		field.String("email").Unique(),
		field.Bool("email_verified"),

		// User real info
		field.String("name"),
		field.String("given_name"),
		field.String("family_name"),
		field.Time("birthdate").Optional().Nillable(),

		// User Info in the service
		field.Text("profile_picture_url").Optional().Nillable(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("hus_sessions", HusSession.Type),
	}
}
