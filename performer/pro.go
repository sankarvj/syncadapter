package performer

import (
	"database/sql"
	"github.com/sankarvj/syncadapter/core"
	"log"
	"reflect"
	"strconv"
)

const (
	NOTHING = 0
	CREATE  = 1
	UPDATE  = 2
)

type Pro struct {
	DBInst          *sql.DB
	Tablename       string
	Localid         int64
	DatabaseChanged bool
}

//Basic
func CreatePro(db *sql.DB) Pro {
	return Pro{db, "", 0, false}
}

//Advanced
func CreateProAdv(db *sql.DB, tablename string, localid int64) Pro {
	return Pro{db, tablename, localid, false}
}

//Using channels prepare cooker for api object and waits for the result.
//Once result receives it updates the database with the server key
func (s *Pro) ApiMeltDown(cooker core.Cooker) chan core.Cooker {
	s.CookForRemote(cooker)
	bridgeForCooker := make(chan core.Cooker)

	go func(bridgeForCooker chan core.Cooker) {
		log.Println("Waiting for API call to finish")
		cooker = <-bridgeForCooker
		log.Println("Channel received successfully ")
		if cooker != nil {
			s.CoolItDown(cooker)
			log.Println("Channel cooled successfully")
		}

	}(bridgeForCooker)

	return bridgeForCooker
}

//Using channels prepare cooker for api object and waits for the result.
//Once result receives it updates the database with the server key
func (s *Pro) ApiDeleteDown(cooker core.Cooker) chan core.Cooker {
	s.CookForRemote(cooker)
	bridgeForCooker := make(chan core.Cooker)

	go func(bridgeForCooker chan core.Cooker) {
		log.Println("Waiting for API call to finish")
		cooker = <-bridgeForCooker
		log.Println("Channel received successfully ")
		if cooker != nil {
			s.DeleteItem(s.Tablename, cooker.LocalId())
			log.Println("Channel cooled successfully")
		} else {
			s.markAsDeletedInLocal()
			log.Println("Channel is still hot")
		}

	}(bridgeForCooker)

	return bridgeForCooker
}

//Basic funcs which sets db value manually

//Update the key and time value of the local db from the server obj (wrapper for outside world)
func (s *Pro) CoolItDown(cooker core.Cooker) {
	log.Println("Cool it down cooker --- ", cooker)
	if cooker != nil {
		s.coolItDown(cooker.LocalId(), cooker.UpdatedAt())
	}
}

//HotId returns serverKey for the localId
func (s *Pro) HotId(tablename string, localid int64) int64 {
	return serverVal(s.DBInst, tablename, strconv.FormatInt(localid, 10))
}

//ColdId returns localId for serverKey
func (s *Pro) ColdId(tableName string, serverId int64) (string, bool) {
	localId, isAvailable := localkey(s.DBInst, tableName, serverId)
	return strconv.FormatInt(localId, 10), isAvailable
}

//Append
func (s *Pro) PrepareLocal(cooker core.Cooker, tablename string) {
	if cooker.LocalId() != 0 {
		cooker.PrepareLocal(true, s.HotId(tablename, cooker.LocalId()))
	} else {
		cooker.PrepareLocal(true, 0)
	}
}

//ColdId returns localId for serverKey
func (s *Pro) DeleteItem(tableName string, serverId int64) bool {
	err := deleteItem(s.DBInst, tableName, serverId)
	if err != nil {
		log.Println("Error deleting item in recorder ", err)
		return false
	}
	return true
}

//Update the key and time value of the local db from the server obj (original implementation)
func (s *Pro) coolItDown(key int64, updated int64) {
	updateKey(s.DBInst, s.Tablename, key, s.Localid, updated)
}

//Update the key and time value of the local db from the server obj (original implementation)
func (s *Pro) markAsDeletedInLocal() {
	markAsDeletedLocally(s.DBInst, s.Tablename, s.Localid)
}

//Logic to find new or updated items from the server list. It returns newitems and updateditems array
func (s *Pro) WhatToDoLogic1(slice interface{}, locallistitems []core.Passer) ([]core.Passer, []core.Passer) {
	serverlistitems := reflect.ValueOf(slice)
	if serverlistitems.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	newItems := make([]core.Passer, 0)
	updatedItems := make([]core.Passer, 0)

	var localitem core.Passer
	for i := 0; i < serverlistitems.Len(); i++ {
		serveritem := serverlistitems.Index(i).Addr().Interface().(core.Passer)
		s.CookFromRemote(serveritem)

		presentInDB := false
		for j := 0; j < len(locallistitems); j++ {
			localitem = locallistitems[j]

			if (serveritem).ServerKey() == localitem.ServerKey() {
				presentInDB = true
				if needUpdate(serveritem.UpdatedAt(), localitem.UpdatedAt()) {
					s.DatabaseChanged = true
					serveritem.(core.Cooker).SetLocalId(localitem.LocalId())
					updatedItems = append(updatedItems, serveritem)
				}
			}
		}
		//Check for new
		if !presentInDB && (serveritem).ServerKey() != 0 { //some rare cases server sends the empty model
			s.DatabaseChanged = true
			newItems = append(newItems, serveritem)
		}
	}

	return newItems, updatedItems
}

// func (s *Pro) WhatToDoLogic2(slice interface{}, locallistitems []core.Cooker) ([]core.Passer, []core.Passer) {
// 	serverlistitems := reflect.ValueOf(slice)
// 	if serverlistitems.Kind() != reflect.Slice {
// 		panic("InterfaceSlice() given a non-slice type")
// 	}

// 	newItems := make([]core.Passer, 0)
// 	updatedItems := make([]core.Passer, 0)

// 	var localitem core.Passer
// 	for i := 0; i < serverlistitems.Len(); i++ {
// 		serveritem := serverlistitems.Index(i).Addr().Interface().(core.Passer)
// 		s.CookFromRemote(serveritem)

// 		presentInDB := false
// 		for j := 0; j < len(locallistitems); j++ {
// 			localitem = locallistitems[j]

// 			if (serveritem).ServerKey() == localitem.ServerKey() {
// 				presentInDB = true
// 				if needUpdate(serveritem.UpdatedAt(), localitem.UpdatedAt()) {
// 					s.DatabaseChanged = true
// 					serveritem.(core.Cooker).SetLocalId(localitem.LocalId())
// 					updatedItems = append(updatedItems, serveritem)
// 				}
// 			}
// 		}
// 		//Check for new
// 		if !presentInDB && (serveritem).ServerKey() != 0 { //some rare cases server sends the empty model
// 			s.DatabaseChanged = true
// 			newItems = append(newItems, serveritem)
// 		}
// 	}

// 	for i := 0; i < len(locallistitems); i++ {
// 		localitem := locallistitems[i]
// 		log.Println("localitem.ServerKey() --- ", localitem.ServerKey())
// 		if localitem.ServerKey() == 0 {
// 			//Try to sync here
// 			periodic := technique.CreatePeriodic(s.DBInst)
// 			periodic.Models = append(periodic.Models, core.Cooker(localitem))
// 			periodic.CheckPeriodic()
// 		}
// 	}

// 	return newItems, updatedItems
// }

func (s *Pro) Push(cooker core.Cooker) bool {
	s.CookForRemote(cooker)

	var remoteUpdated bool
	if cooker.ServerKey() != 0 { //update
		remoteUpdated = cooker.Signal(core.BASIC_UPDATE)
	} else { //create
		remoteUpdated = cooker.Signal(core.BASIC_CREATE)
	}

	if remoteUpdated {
		//Using LocalId here is very misleading. Also we can't sure that the user implementation always update ID as serverkey
		s.coolItDown(cooker.LocalId(), cooker.UpdatedAt())
	}

	return remoteUpdated
}
