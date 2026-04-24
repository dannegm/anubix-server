package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type ShareToken struct {
	ent.Schema
}

func (ShareToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.Enum("type").Values("permanent", "expiring", "one_time"),
		field.String("api_key").Optional().Nillable(),
		field.String("api_secret_hash").Optional().Nillable().Sensitive(),
		field.Text("ciphertext").Optional().Nillable(),
		field.String("iv").Optional().Nillable(),
		field.String("auth_tag").Optional().Nillable(),
		field.Int("use_count").Default(0),
		field.Time("used_at").Optional().Nillable(),
		field.Time("expires_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (ShareToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("entry", Entry.Type).Ref("share_tokens").Unique().Required(),
		edge.From("created_by", User.Type).Ref("share_tokens").Unique().Required(),
	}
}
