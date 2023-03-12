package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Subdomain holds the schema definition for the Subdomain entity.
type Subdomain struct {
	ent.Schema
}

// Fields of the Subdomain.
func (Subdomain) Fields() []ent.Field {
	return []ent.Field{
		field.Int("service_id"),
		field.String("subdomain"),
		field.String("url"),
	}
}

// Edges of the Subdomain.
func (Subdomain) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).Ref("subdomains").Unique().Field("service_id").Required(),
	}
}
