package sqlite

import (
	"time"
)

type NetworkMgr struct {
	Id         int64
	Subnet     string
	BindIp     string
	CreateTime time.Time
	UpdateTime time.Time
}

// insert bridge subnet
func (network *NetworkMgr) Insert(u NetworkMgr) error {
	db = GetDb()
	sql := `insert into network_mgr (subnet, bind_ip) values(?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Subnet, u.BindIp)
	return err
}

// query bridge subnet ip,for new assigned ip
func (network *NetworkMgr) QueryBySubnet(bridgeSubnet string) (l []*NetworkMgr, e error) {
	db = GetDb()
	sql := `select * from network_mgr where subnet=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(bridgeSubnet)
	if err != nil {
		return nil, err
	}
	var result = make([]*NetworkMgr, 0)
	for rows.Next() {
		var subnet, bindIp string
		var id int64
		var createTime, updateTime time.Time
		rows.Scan(&id, &subnet, &bindIp, &createTime, &updateTime)
		result = append(result, &NetworkMgr{id, subnet, bindIp, createTime, updateTime})
	}
	return result, nil
}

func (network *NetworkMgr) DeleteByBindIp(bindIp string) (bool, error) {
	db = GetDb()
	sql := `delete from network_mgr where bind_ip=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(bindIp)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}
