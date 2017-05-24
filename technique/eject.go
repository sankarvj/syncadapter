package technique

import (
	"database/sql"
)

type Erasesync struct {
	DBInst     *sql.DB
	Tablenames []string
}

//When you completly screwed up
func DropDB() {

}

//When you partially screwed up
func DeleteTable(tablename string) {

}

//When the data stored in frozen state for a very long even after many attempt to sync that data
func EliminateRottonData(tablename string) {

}
