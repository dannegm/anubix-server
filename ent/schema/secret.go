package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type Secret struct {
	ent.Schema
}

func (Secret) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.Text("ciphertext").NotEmpty(),
		field.String("iv").NotEmpty(),
		field.String("auth_tag").NotEmpty(),
	}
}

func (Secret) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("block", Block.Type).Ref("secrets").Unique().Required(),
	}
}
