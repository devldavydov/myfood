package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Bundle holds the schema definition for the Bundle entity.
type Bundle struct {
	ent.Schema
}

// Fields of the Bundle.
func (Bundle) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("userid"),
		field.String("key"),
		field.JSON("data", map[string]float64{}),
	}
}

// Edges of the Bundle.
func (Bundle) Edges() []ent.Edge {
	return nil
}

// Indexes of the Bundle
func (Bundle) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Fields("userid", "key").
			Unique(),
	}
}
