package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Journal holds the schema definition for the Journal entity.
type Journal struct {
	ent.Schema
}

// Fields of the Journal.
func (Journal) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("userid"),
		field.Time("timestamp"),
		field.Int64("meal"),
		field.Float("foodweight"),
	}
}

// Edges of the Journal.
func (Journal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.
			From("food", Food.Type).
			Ref("journals").
			Unique().
			Required(),
	}
}

// Indexes of the Journal
func (Journal) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Fields("userid", "timestamp", "meal").
			Edges("food").
			Unique(),
	}
}
