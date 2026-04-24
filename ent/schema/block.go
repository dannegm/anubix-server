package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type Block struct {
	ent.Schema
}

func (Block) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.String("label").NotEmpty(),
		field.Int("sort_order").Default(0),
	}
}

func (Block) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("entry", Entry.Type).Ref("blocks").Unique().Required(),
		edge.To("secrets", Secret.Type),
	}
}
