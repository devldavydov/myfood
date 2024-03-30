package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Weight holds the schema definition for the Weight entity.
type Weight struct {
	ent.Schema
}

// Fields of the Weight.
func (Weight) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("userid"),
		field.Time("timestamp"),
		field.Float("value"),
	}
}

// Edges of the Weight.
func (Weight) Edges() []ent.Edge {
	return nil
}

// Indexes of the Weight
func (Weight) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("userid", "timestamp").Unique(),
	}
}
