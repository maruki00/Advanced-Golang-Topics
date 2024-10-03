package main

type OrderEvent interface {
	Event
	OrderId() int
}

type OrderDispatched struct {
	orderId int
}




