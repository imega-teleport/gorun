package teleport

import (
	"fmt"
	"time"

	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

type TeleportItem struct {
	GUID UUID
	Type string
	ID   int
	Date time.Time
}

func (t TeleportItem) SizeOf() int {
	return len(t.GUID) + len(t.Type) + lengthDefineIndex + lengthDefineDate
}

func (w *Wpwc) BuilderTeleportItem() builder {
	return builder{
		squirrel.Insert(fmt.Sprintf("%steleport_item", w.Prefix)).Columns("guid", "type", "id", "date"),
	}
}

func (b *builder) AddTeleportItem(i TeleportItem) {
	*b = builder{
		b.Values(
			i.GUID,
			i.Type,
			squirrel.Expr(fmt.Sprintf("@max_post_id+%s", i.GUID.ToVar())),
			i.Date.Format("2006-01-02 15:04:05"),
		),
	}
}
