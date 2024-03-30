package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Food holds the schema definition for the Food entity.
type Food struct {
	ent.Schema
}

// Fields of the Food.
func (Food) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").Unique(),
		field.String("name").NotEmpty(),
		field.String("brand").Optional(),
		field.Float("cal100"),
		field.Float("prot100"),
		field.Float("fat100"),
		field.Float("carb100"),
		field.String("comment").Optional(),
	}
}

// Edges of the Food.
func (Food) Edges() []ent.Edge {
	return []ent.Edge{
		edge.
			To("journals", Journal.Type).
			Annotations(entsql.OnDelete(entsql.Restrict)),
	}
}
