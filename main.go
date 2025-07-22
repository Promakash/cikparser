package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Поменять имена файлов на свои если хочется.
	const (
		InputXMLFilename  = "report.xml"
		OutputCSVFilename = "report.csv"
	)

	report, err := ParseXMLFile(InputXMLFilename)
	if err != nil {
		fmt.Println("error while parsing xml file:", err)
		return
	}

	err = CreateParsedTable(report, OutputCSVFilename)
	if err != nil {
		fmt.Println("error while creating parsed table:", err)
		return
	}
}

func ParseXMLFile(filename string) (*CIKReport, error) {
	// Читаем файл.
	fileRaw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Добавляем поля, чтобы утилита успешно распарсила файл.
	fileWithAdditionalFields := "<Persons>" + string(fileRaw) + "</Persons>"
	// Переводим в байты чтобы анмаршаллер смог распарсить.
	fileWithAdditionalFieldsByteForm := []byte(fileWithAdditionalFields)

	// Анмаршаллим в структуру.
	var report CIKReport
	err = xml.Unmarshal(fileWithAdditionalFieldsByteForm, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func CreateParsedTable(report *CIKReport, outputFileName string) error {
	// Создаем CSV файл.
	csvFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = csvFile.Close()
	}()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Записываем заголовки csv файла.
	headers := []string{"Фамилия", "Имя", "Отчество", "ДатаРожд", "Серия паспорта", "Номер паспорта"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Записываем данные в csv файл.
	for _, persona := range report.Personas {
		// В ответах серия документа разделена пробелами, нормализуем их.
		correctedDocSeries := strings.ReplaceAll(persona.Info.Doc.Series, " ", "")

		record := []string{
			persona.Info.FIO.LastName,
			persona.Info.FIO.FirstName,
			persona.Info.FIO.Patronymic,
			persona.Info.FIO.Birthday,
			correctedDocSeries,
			persona.Info.Doc.Number,
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	fmt.Println("CSV файл успешно создан: ", outputFileName)
	return nil
}
