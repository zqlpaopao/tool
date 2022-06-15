package mysql

import (
	"strings"
)

const (
	UPDATE = "UPDATE `"
	SET    = "` SET `"
	CASE   = "` = CASE `"
	WHEN   = " WHEN '"
	THEN   = "' THEN '"
	END    = " END ,`"
	WHERE  = " WHERE `"
	IN     = "` IN ("
)

//UpdateInfo Update information to batch
type UpdateInfo struct {
	//TableName table name
	TableName string

	//SET update_column = CASE where_column
	SetCase map[string]string

	// WHEN where_column == 1 THEN update_column = 3
	WhenThen map[string]map[string]string

	//WHERE where_column IN (where_column_value1,where_column_value1,where_column_value1)
	Where map[string][]string
}

//GetUpdateBatch Batch obtain and modify SQL statements of different fields in different tables
//请求参数
/*
	item := mysql.UpdateInfo{
			//表名
			TableName: "table_name",

			SetCase: map[string]string{
				"name" : "id",//要修改的列名---where条件的列名
				"addr" : "id",//要修改的列名---where条件的列名
			},
			WhenThen: map[string]map[string]string{
				"name" :{
					"1" : "name1",//当id =1时候，name的值修改为name1
					"2" : "name12",//当id =2时候，name的值修改为name2
					"3" : "name13",//当id =3时候，name的值修改为name3
				},
				"addr":{
					"1" : "addr1",//当id =1时候，addr的值修改为addr1
					"2" : "addr2",//当id =2时候，addr的值修改为addr2
					"3" : "addr3",//当id =3时候，addr的值修改为addr3
				},
			},
			Where: map[string][]string{
				"id":{//where条件的列
					"1",//上面的值
					"2",
					"3",
				},
			},
		}
*/
/*
	[
		UPDATE categories
		SET display_order = CASE id
		  WHEN 1 THEN 3
		  WHEN 2 THEN 4
		  WHEN 3 THEN 5
		END,
		title = CASE id
		  WHEN 1 THEN 'New Title 1'
		  WHEN 2 THEN 'New Title 2'
		  WHEN 3 THEN 'New Title 3'
		END
		WHERE id IN (1,2,3)

		UPDATE categories
		SET display_order = CASE id
		  WHEN 1 THEN 3
		  WHEN 2 THEN 4
		  WHEN 3 THEN 5
		END,
		title = CASE id
		  WHEN 1 THEN 'New Title 1'
		  WHEN 2 THEN 'New Title 2'
		  WHEN 3 THEN 'New Title 3'
		END
		WHERE id IN (1,2,3)
	]
*/
func GetUpdateBatch(info []*UpdateInfo) (updateSql []string) {
	updateSql = make([]string, 0, len(info))
	for i := 0; i < len(info); i++ {
		updateSql = append(updateSql, GetUpdateBatchItem(info[i]))
	}
	return
}

//GetUpdateBatchItem SQL statements for different fields of batch tables
/*
	UPDATE categories
	SET display_order = CASE id
	  WHEN 1 THEN 3
	  WHEN 2 THEN 4
	  WHEN 3 THEN 5
	END,
	title = CASE id
	  WHEN 1 THEN 'New Title 1'
	  WHEN 2 THEN 'New Title 2'
	  WHEN 3 THEN 'New Title 3'
	END
	WHERE id IN (1,2,3)
*/
func GetUpdateBatchItem(info *UpdateInfo) (sql string) {
	sql = UPDATE + info.TableName + SET
	for k, v := range info.SetCase {
		sql += k + CASE + v + "`" + GetWhenThen(info.WhenThen[k])
	}
	return strings.TrimRight(sql, ",`") + GetWhere(info.Where)
}

//GetWhenThen Get when case
/*
	WHEN 1 THEN 3
  	WHEN 2 THEN 4
  	WHEN 3 THEN 5
	END,
*/
func GetWhenThen(info map[string]string) (sql string) {
	for k, v := range info {
		sql += WHEN + k + THEN + v + "' "
	}
	return sql + END
}

//GetWhere
//WHERE id IN (1,2,3)
func GetWhere(info map[string][]string) (where string) {
	where = WHERE
	for k, v := range info {
		where += k + IN
		for i := 0; i < len(v); i++ {
			where += "'" + v[i] + "',"
		}
	}
	return strings.TrimRight(where, ",") + ") "
}
