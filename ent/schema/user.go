package schema

import (
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
		field.UUID("uuid", uuid.UUID{}).Default(uuid.New),
		field.String("google_sub").Unique(),

		field.String("email").Unique(),
		field.Bool("email_verified"),

		// User real info
		field.String("name"),
		field.Time("birthday"),
		field.String("given_name"),
		field.String("family_name"),

		// User Info in the service
		field.Text("google_profile_picture"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("groups", Group.Type).Ref("users"),
	}
}
