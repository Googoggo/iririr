package datamodels

type Product struct {
	ID int64 `json:"id" sql:"ID" database2:"ID"`
	ProductName string `json:"ProductName" sql:"productName" database2:"ProductName"`
	ProductNum int64 `json:"ProductNum" sql:"productNum" database2:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" database2:"ProductImage"`
	ProductUrl string `json:"ProductUrl" sql:"productUrl" database2:"ProductUrl"`
}
