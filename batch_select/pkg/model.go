package pkg

type MaxMinInfo struct {
	Max string `gorm:"column:ma"`
	Min string `gorm:"column:mi"`
}

// getMaxMinInfo Get the maximum and minimum values that meet the query criteria
func (o *Option[T]) getMaxMinInfo() (maxMinData MaxMinInfo, err error) {
	err =
		o.mysqlCli.
			Table(o.table).
			Select("min("+o.orderColumn+") as mi,"+"max("+o.orderColumn+") as ma").
			Where(o.sqlWhere, o.whereCase...).
			Take(&maxMinData).
			Error
	return
}

// getResInfo Get data according to conditions
func (o *Option[T]) getResInfo(sqlWhere string, whereCase []interface{}) (res []T, err error) {
	res = make([]T, 0, o.limit)
	err = o.mysqlCli.
		Table(o.table).
		Select(o.selectFiled).
		Where(sqlWhere, whereCase...).
		//Limit(o.limit).
		//Order(o.orderColumn).
		Find(&res).
		Error
	return
}
