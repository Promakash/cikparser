package main

import "encoding/xml"

type CIKReport struct {
	XMLName  xml.Name  `xml:"Persons"`
	Personas []Persona `xml:"Persona"`
}

type Persona struct {
	Info PersonalInfo `xml:"ПерсИнфо"`
}

type PersonalInfo struct {
	FIO FIOD     `xml:"ФИОД"`
	Doc Document `xml:"Документ"`
}

type FIOD struct {
	LastName   string `xml:"Фамилия,attr"`
	FirstName  string `xml:"Имя,attr"`
	Patronymic string `xml:"Отчество,attr"`
	Birthday   string `xml:"ДатаРожд,attr"`
}

type Document struct {
	Series string `xml:"Серия,attr"`
	Number string `xml:"Номер,attr"`
}
