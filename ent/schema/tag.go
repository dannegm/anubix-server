package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type Tag struct {
	ent.Schema
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.String("name").NotEmpty(),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("tags").Unique().Required(),
		edge.To("entry_tags", EntryTag.Type),
	}
}
