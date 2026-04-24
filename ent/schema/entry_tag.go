package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

type EntryTag struct {
	ent.Schema
}

func (EntryTag) Fields() []ent.Field {
	return []ent.Field{}
}

func (EntryTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("entry", Entry.Type).Ref("entry_tags").Unique().Required(),
		edge.From("tag", Tag.Type).Ref("entry_tags").Unique().Required(),
	}
}
