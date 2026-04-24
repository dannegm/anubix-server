package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/oklog/ulid/v2"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return ulid.Make().String()
			}).
			Immutable().
			Unique(),
		field.String("email").Unique().NotEmpty(),
		field.String("auth_hash").Sensitive(),
		field.String("salt").NotEmpty(),
		field.Time("email_verified_at").Optional().Nillable(),
		field.String("two_factor_secret").Optional().Nillable().Sensitive(),
		field.Bool("two_factor_enabled").Default(false),
		field.String("password_reset_token").Optional().Nillable(),
		field.Time("password_reset_expires_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vaults", Vault.Type),
		edge.To("devices", Device.Type),
		edge.To("sessions", Session.Type),
		edge.To("tags", Tag.Type),
		edge.To("audit_logs", AuditLog.Type),
		edge.To("share_tokens", ShareToken.Type),
	}
}
