package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// UserSettings holds the schema definition for the UserSettings entity.
type UserSettings struct {
	ent.Schema
}

// Fields of the UserSettings.
func (UserSettings) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("userid").Unique(),
		field.Float("cal_limit"),
	}
}

// Edges of the UserSettings.
func (UserSettings) Edges() []ent.Edge {
	return nil
}
