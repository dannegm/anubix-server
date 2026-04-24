package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type Entry struct {
	ent.Schema
}

func (Entry) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.String("label").NotEmpty(),
		field.JSON("icon", map[string]interface{}{}).Optional(),
		field.String("preview").Optional(),
		field.Bool("has_otp").Default(false),
		field.Bool("is_favorite").Default(false),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Entry) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("vault", Vault.Type).Ref("entries").Unique().Required(),
		edge.To("blocks", Block.Type),
		edge.To("attachments", Attachment.Type),
		edge.To("share_tokens", ShareToken.Type),
		edge.To("audit_logs", AuditLog.Type),
		edge.To("entry_tags", EntryTag.Type),
	}
}
