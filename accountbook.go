package main

import (
	"errors"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
)

// TODO:
// 使用命令为频道建立账本
// 频道成员可进行记账
// 若频道不可访问，则删除对应账本

var db *scribble.Driver

var cacheRecord []accountRecord

type moneyRecord struct {
	User    string
	Time    int64
	Money   float64
	Comment string
}

type userRecord struct {
	User  string
	Money float64
}

type accountRecord struct {
	Id       string
	URecords []userRecord
	MRecords []moneyRecord
}

func (a *accountRecord) Add(user string, money float64, comment string) error {
	a.MRecords = append(a.MRecords, moneyRecord{user, time.Now().Unix(), money, comment})
	var found bool = false
	for i, v := range a.URecords {
		if v.User == user {
			a.URecords[i].Money += money
			found = true
		}
	}
	if !found {
		a.URecords = append(a.URecords, userRecord{user, money})
	}
	if err := db.Write("db", "allAccountRecord", cacheRecord); err != nil {
		return errors.New("账本保存失败")
	}
	return nil
}

func accountBookInit() {
	db, _ = scribble.New("./database", nil)
	db.Read("db", "allAccountRecord", &cacheRecord)
	// TODO: 核算余额
	// for i,v:=range cacheRecord{}
	// fmt.Printf("%v\r\n", cacheRecord)
}

func accountBookCreate(id string) error {
	for _, v := range cacheRecord {
		if v.Id == id {
			return errors.New("账本已经存在")
		}
	}
	cacheRecord = append(cacheRecord, accountRecord{id, []userRecord{}, []moneyRecord{}})
	if err := db.Write("db", "allAccountRecord", cacheRecord); err != nil {
		return errors.New("账本保存失败")
	}
	return nil
}

func accountBookGetSummary(id string) ([]userRecord, error) {
	for _, v := range cacheRecord {
		if v.Id == id {
			return v.URecords, nil
		}
	}
	return nil, errors.New("未找到账本")
}

func accountBookRecordAdd(groupId string, user string, money float64, comment string) error {
	var found bool = false
	for i, v := range cacheRecord {
		if v.Id == groupId {
			cacheRecord[i].Add(user, money, comment)
			found = true
			break
		}
	}
	if !found {
		return errors.New("没有注册账本")
	}

	if err := db.Write("db", "allAccountRecord", cacheRecord); err != nil {
		return errors.New("账本保存失败")
	}
	// fmt.Printf("%v\r\n", cacheRecord)
	return nil
}