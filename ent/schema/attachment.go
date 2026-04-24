package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type Attachment struct {
	ent.Schema
}

func (Attachment) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.String("filename").NotEmpty(),
		field.String("mime_type").NotEmpty(),
		field.Int("size_bytes"),
		field.Bytes("ciphertext").NotEmpty(),
		field.String("iv").NotEmpty(),
		field.String("auth_tag").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Attachment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("entry", Entry.Type).Ref("attachments").Unique().Required(),
	}
}
