/**
* @Author: yanKoo
* @Date: 2019/3/18 11:34
* @Description:
 */
package dbops

import "log"

func SelectDeviceIdsByGroupId(gid int) ([]int, error) {
	stmtOut, err := dbConn.Prepare("SELECT device_id FROM group_device WHERE group_id = ?")
	if err != nil {
		return nil, err
	}

	var res []int
	rows, err := stmtOut.Query(gid)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return res, err
		}
		res = append(res, id)
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("statement close fail.")
		}
	}()
	return res, nil
}
