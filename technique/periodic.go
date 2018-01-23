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
	"github.com/sankarvj/syncadapter/utils"
	"log"
)

type Periodic struct {
	DBInst *sql.DB
	Models []core.Cooker
}

func CreatePeriodic(db *sql.DB) Periodic {
	return Periodic{db, make([]core.Cooker, 0)}
}

func SingleKickPeriodic(db *sql.DB, model core.Cooker) Periodic {
	periodic := Periodic{db, make([]core.Cooker, 0)}
	periodic.Models = append(periodic.Models, model)
	periodic.CheckPeriodic()
	return periodic
}

func (g *Periodic) CheckPeriodic() {
	var basemodel core.BaseModel
	var technique core.Technique

	for i := 0; i < len(g.Models); i++ {
		basemodels := scanFrozenData(g.DBInst, utils.Tablename(g.Models[i]))
		for j := 0; j < len(basemodels); j++ {
			basemodel = basemodels[j]
			if basemodel.Key == nil { //Create
				technique = core.BASIC_CREATE
			} else if basemodel.Updated == -1 && basemodel.Synced == false { // record deleted
				technique = core.BASIC_DELETE
			} else { //Update
				technique = core.BASIC_UPDATE
			}
			//Fire signal
			g.Models[i].SetLocalId(basemodel.Id)
			g.Models[i].Signal(technique)
		}
	}
}

func scanFrozenData(db *sql.DB, tablename string) []core.BaseModel {
	sql_readall := `
	SELECT Id,Key,Updated,Synced FROM ` + tablename + `
	WHERE Synced = 0
	`
	basemodels := make([]core.BaseModel, 0)
	rows, err := db.Query(sql_readall)
	defer closeRows(rows)
	if err != nil {
		log.Println("Error reading scanFrozenData ", err)
		return basemodels
	}

	var basemodel core.BaseModel
	for rows.Next() {
		err = rows.Scan(&basemodel.Id, &basemodel.Key, &basemodel.Updated, &basemodel.Synced)
		if err != nil {
			log.Println("Error scan ", err)
			return basemodels
		}
		basemodels = append(basemodels, basemodel)
	}

	return basemodels
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}
