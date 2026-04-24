package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type AuditLog struct {
	ent.Schema
}

func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.Enum("action").Values(
			"login",
			"logout",
			"vault_created",
			"vault_deleted",
			"entry_created",
			"entry_updated",
			"entry_deleted",
			"entry_viewed",
			"entry_shared",
			"password_changed",
		),
		field.String("ip_address").Optional().Nillable(),
		field.String("user_agent").Optional().Nillable(),
		field.JSON("metadata", map[string]interface{}{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("audit_logs").Unique().Required(),
		edge.From("device", Device.Type).Ref("audit_logs").Unique(),
		edge.From("entry", Entry.Type).Ref("audit_logs").Unique(),
	}
}
