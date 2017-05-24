package technique

//#Core goals
//* pull changes in the background without affecting the user experiance

//#Needs server side implementation
//* no

//#Logic
//* init time_specific or network_specific

//#TODO
//* find a way to sense network reconnectivity
//* give 2 minutes buffer in the event of network reconnectivity. This will help reduce server load.

import (
	"database/sql"
	"github.com/sankarvj/syncadapter/core"
	"github.com/sankarvj/syncadapter/performer"
)

type Periodic struct {
	DBInst *sql.DB
	Models []core.Cooker
}

func CreatePeriodic(db *sql.DB) Periodic {
	return Periodic{db, make([]core.Cooker, 0)}
}

func (g *Periodic) CheckPeriodic() {
	var basemodel core.BaseModel
	var technique core.Technique

	for i := 0; i < len(g.Models); i++ {
		basemodels := performer.ScanFrozenData(g.DBInst, performer.Tablename(g.Models[i]))
		for j := 0; j < len(basemodels); j++ {
			basemodel = basemodels[j]
			if basemodel.Key == 0 { //Create
				technique = core.BASIC_CREATE
			} else { //Update
				technique = core.BASIC_UPDATE
			}
			//Fire signal
			g.Models[i].SetLocalId(basemodel.Id)
			g.Models[i].Signal(technique)
		}
	}
}
