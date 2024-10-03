package main

type Event interface {
	Name() string
}

type GeneralError string

func (e GeneralError) Name() string {
	return "event.general.error"
}
