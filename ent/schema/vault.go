package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type Vault struct {
	ent.Schema
}

func (Vault) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.String("name").NotEmpty(),
		field.Text("encrypted_vault_key").NotEmpty(),
		field.String("vault_key_iv").NotEmpty(),
		field.String("vault_key_auth_tag").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Vault) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("vaults").Unique().Required(),
		edge.To("entries", Entry.Type),
	}
}
