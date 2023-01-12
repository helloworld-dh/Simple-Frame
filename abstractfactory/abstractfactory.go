package abstractfactory

import "fmt"

// OrderMainDao 订单的主记录
type OrderMainDao interface {
	SaveOrderMain()
}

// OrderDetailDao 订单的详情记录
type OrderDetailDao interface {
	SaveOrderDetail()
}

// DaoFactory 抽象工厂接口
type DaoFactory interface {
	CreateOrderMainDao() OrderMainDao
	CreateOrderDetailDao() OrderDetailDao
}

type RDBMainDao struct {
}

func (R *RDBMainDao) SaveOrderMain() {
	fmt.Print("rdb main save\n")
}

type RDBDetailDao struct {
}

func (R RDBDetailDao) SaveOrderDetail() {
	fmt.Print("rdb detail save\n")
}

type RDBDAOFactory struct {
}

func (R *RDBDAOFactory) CreateOrderMainDao() OrderMainDao {
	return &RDBMainDao{}
}

func (R *RDBDAOFactory) CreateOrderDetailDao() OrderDetailDao {
	return &RDBDetailDao{}
}

type XMLMainDao struct {
}

func (X *XMLMainDao) SaveOrderMain() {
	fmt.Print("xml main save\n")
}

type XMLDetailDao struct {
}

func (X *XMLDetailDao) SaveOrderDetail() {
	fmt.Print("xml detail save\n")
}

type XMLDaoFactory struct {
}

func (X *XMLDaoFactory) CreateOrderMainDao() OrderMainDao {
	return &XMLMainDao{}
}

func (X *XMLDaoFactory) CreateOrderDetailDao() OrderDetailDao {
	return &XMLDetailDao{}
}
