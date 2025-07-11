package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Поменять имена файлов на свои если хочется
	const (
		InputXMLFilename  = "report.xml"
		OutputCSVFilename = "report.csv"
	)

	report, err := ParseXMLFile(InputXMLFilename)
	if err != nil {
		return
	}

	err = CreateParsedTable(report, OutputCSVFilename)
	if err != nil {
		return
	}
}

func ParseXMLFile(filename string) (*CIKReport, error) {
	// читаем файл
	fileRaw, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// Добавляем поля, чтобы утилита успешно распарсила файл
	fileWithAdditionalFields := "<Persons>" + string(fileRaw) + "</Persons>"
	// Переводим в байты чтобы анмаршаллер смог распарсить
	fileWithAdditionalFieldsByteForm := []byte(fileWithAdditionalFields)

	var report CIKReport
	err = xml.Unmarshal(fileWithAdditionalFieldsByteForm, &report)
	if err != nil {
		fmt.Println("Error unmarshalling report:", err)
		return nil, err
	}

	return &report, nil
}

func CreateParsedTable(report *CIKReport, outputFileName string) error {
	// Создаем CSV файл
	csvFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return err
	}
	defer func() {
		_ = csvFile.Close()
	}()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Записываем заголовки
	headers := []string{"Фамилия", "Имя", "Отчество", "Серия паспорта", "Номер паспорта"}
	if err := writer.Write(headers); err != nil {
		fmt.Println("Error writing headers:", err)
		return err
	}

	// Записываем данные
	for _, persona := range report.Personas {
		// В ответах серия документа разделена пробелами
		correctedDocSeries := strings.ReplaceAll(persona.Info.Doc.Series, " ", "")

		record := []string{
			persona.Info.FIO.LastName,
			persona.Info.FIO.FirstName,
			persona.Info.FIO.Patronymic,
			correctedDocSeries,
			persona.Info.Doc.Number,
		}

		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record:", err)
			return err
		}
	}

	fmt.Println("CSV файл успешно создан: report.csv")
	return nil
}
