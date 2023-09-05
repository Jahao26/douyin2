package repository

import (
	"fmt"
	"strconv"
	"time"
)

func GetMassage(uid int64, to_uid int64, massageList *[]*Massage) error {
	var timestamp int64
	UserId := fmt.Sprintf("user:%d", uid)

	//pipe := rdb3.TxPipeline()
	exist, _ := rdb3.Exists(c, UserId).Result()
	if exist != 1 {
		//fmt.Println("reset time:", timestamp)
		rdb3.Set(c, UserId, 0, 1*time.Minute) //设置聊天时间戳
	}
	s, _ := rdb3.Get(c, UserId).Result()
	timestamp, _ = strconv.ParseInt(s, 10, 64)
	//fmt.Println("before time:", timestamp)
	timestamp, err := NewMassageDAO().QuaryMassage(uid, to_uid, timestamp, massageList)
	//fmt.Println("after time:", timestamp)
	if err != nil {
		return err
	}
	rdb3.Del(c, UserId)
	rdb3.Set(c, UserId, timestamp, 1*time.Minute)
	//_, err = pipe.Exec(c)
	return err
}

func GetLasttime(uid int64) (int64, error) {
	var lasttime int64
	Usertime := fmt.Sprintf("usertime:%d", uid)
	exist, _ := rdb3.Exists(c, Usertime).Result()
	if exist != 1 {
		rdb3.Set(c, Usertime, time.Now().Unix(), 1*time.Minute) //设置聊天时间戳
		return time.Now().Unix(), nil
	}
	s, _ := rdb3.Get(c, Usertime).Result()
	lasttime, _ = strconv.ParseInt(s, 10, 64)
	return lasttime, nil
}
