package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type Device struct {
	ent.Schema
}

func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.String("name").NotEmpty(),
		field.String("fingerprint").Unique().NotEmpty(),
		field.Enum("device_type").Values("web", "ios", "android", "desktop"),
		field.Time("last_seen_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("devices").Unique().Required(),
		edge.To("sessions", Session.Type),
		edge.To("audit_logs", AuditLog.Type),
	}
}
