package core

import (
	"time"
)

type BaseModel struct {
	Id      int64 //local id
	Key     int64 //server id
	Updated int64 //last updated time
	Synced  bool  //synced or not
}

//Cooker interface
type Cooker interface {
	UpdatedAt() int64
	LocalId() int64
	ServerKey() int64
	SetLocalId(id int64)
	SetServerKey(key int64)
	PrepareLocal(forced bool, key int64)
	Signal(technique Technique) bool
}

//Cooker implementations
// set forced = true if preparing object before update
func (basemodel *BaseModel) PrepareLocal(forced bool, key int64) {
	basemodel.Key = key
	if basemodel.Id == 0 || forced { //storing ticket originally created at client
		basemodel.Synced = false
		basemodel.Updated = currentTime()
	} else { //storing ticket originally created at server
		basemodel.Synced = true
	}
}

func (basemodel *BaseModel) SetLocalId(id int64) {
	basemodel.Id = id
}

func (basemodel *BaseModel) SetServerKey(key int64) {
	basemodel.Key = key
}

func (basemodel BaseModel) ServerKey() int64 {
	return basemodel.Key
}

func (basemodel BaseModel) UpdatedAt() int64 {
	return basemodel.Updated
}

func (basemodel BaseModel) LocalId() int64 {
	return basemodel.Id
}

//Passer interface
type Passer interface {
	ServerKey() int64
	UpdatedAt() int64
	LocalId() int64
}

type Technique int64

const (
	BASIC_CREATE Technique = iota
	BASIC_UPDATE
	BASIC_NEW_ITEM
	BASIC_UPDATED_ITEM
	BASIC_DELETED_ITEM
)

func currentTime() int64 {
	return milliSeconds(time.Now())
}

func milliSeconds(now time.Time) int64 {
	return now.UTC().Unix() * 1000
}
