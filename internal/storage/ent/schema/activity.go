package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Activity holds the schema definition for the Activity entity.
type Activity struct {
	ent.Schema
}

// Fields of the Activity.
func (Activity) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("userid"),
		field.Time("timestamp"),
		field.Float("active_cal"),
	}
}

// Edges of the Activity.
func (Activity) Edges() []ent.Edge {
	return nil
}

// Indexes of the Activity
func (Activity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("userid", "timestamp").Unique(),
	}
}
