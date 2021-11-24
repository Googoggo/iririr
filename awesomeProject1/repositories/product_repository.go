package repositories

import (
	"awesomeProject1/common"
	"awesomeProject1/datamodels"
	"database/sql"
	"strconv"
)

//创建接口 ， 实现接口

type IProduct interface{
	//链接数据库
	Conn()(error)
	Insert(product * datamodels.Product)(int64, error)
	Delete(int64) bool
	Update(product * datamodels.Product) error
	SelectByKey(int64)(*datamodels.Product, error)
	SelectAll()([]* datamodels.Product, error)
	SubProductNum(productID int64)error
}

type ProductManager struct {
	table string
	mysqlConn *sql.DB
}

func NewProductManager(table string,db *sql.DB) IProduct{
	return &ProductManager{table:table,mysqlConn:db}
}

func (p *ProductManager) Conn()(err error){
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == ""{
		p.table = "product"
	}
	return
}

func (p *ProductManager) Insert(product *datamodels.Product)(productId int64,err error){
	//1.判断链接是否存在
	if err = p.Conn();err != nil{
		return
	}
	//2.准备sql
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	//3.传入参数
	result ,errStmt := stmt.Exec(product.ProductName,product.ProductNum,
		product.ProductImage,product.ProductUrl)
	if errStmt != nil{
		return 0,errStmt
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(productID int64) bool{
	//判断是否存在
	if err := p.Conn();err != nil{
		return false
	}
	sql := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err == nil {
		return false
	}
	_, err = stmt.Exec(productID)
	if err != nil{
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product)(err error){
	if err := p.Conn(); err != nil{
		return nil
	}
	sql := "Update product set productName=?,productNum=?,productImage=?" +
		"productUrl=? where ID=" +strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil{
		return
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage,
		product.ProductUrl)
	if err != nil{
		return
	}
	return
}

func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product,err error){
	//1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where ID =" +strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return

}

func (p *ProductManager) SelectAll()(productArray [] *datamodels.Product, err error){
	if err = p.Conn();err!= nil {
		return nil, err
	}
	sql := "Select *from "+ p.table
	rows, errRow := p.mysqlConn.Query(sql)
	defer rows.Close()
	if errRow != nil {
		return nil, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0{
		return nil, nil
	}
	for _,v := range result{
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return productArray, err
}

func (p *ProductManager) SubProductNum(productID int64) error{
	if err := p.Conn(); err!= nil{
		return err
	}
	sql := "update "+p.table+" set "+" productNum=productNum-1 where ID="+strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil{
		return err
	}
	_,err = stmt.Exec()
	return err
}