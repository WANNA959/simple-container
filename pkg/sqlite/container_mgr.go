package sqlite

import "time"

type ContainerMgr struct {
	Id         int64
	Pid        string
	Veth       string
	CreateTime time.Time
	UpdateTime time.Time
}

// insert pid-veth
func (network *ContainerMgr) Insert(u ContainerMgr) error {
	db = GetDb()
	sql := `insert into container_mgr (pid, veth) values(?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Pid, u.Veth)
	return err
}

// for veth clean
func (network *ContainerMgr) QueryByPid(cpid string) (l []*ContainerMgr, e error) {
	db = GetDb()
	sql := `select * from container_mgr where pid=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(cpid)
	if err != nil {
		return nil, err
	}
	var result = make([]*ContainerMgr, 0)
	for rows.Next() {
		var pid, veth string
		var id int64
		var createTime, updateTime time.Time
		rows.Scan(&id, &pid, &veth, &createTime, &updateTime)
		result = append(result, &ContainerMgr{id, pid, veth, createTime, updateTime})
	}
	return result, nil
}
