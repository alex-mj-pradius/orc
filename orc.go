package orc

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/godror/godror"
)

// Oracle - connect to data base
var Oracle OracleConnect

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////// ORACLE /////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// OracleConnect - ...data base
type OracleConnect struct {
	DBHost  string `xml:"db-host"`  // Oracle.dbHost = "192.168.XXX.XXX:1521/AAAA"
	DBLogin string `xml:"db-login"` // Oracle.dbLogin = "oracle-user"
	DBPass  string `xml:"db-pass"`  // Oracle.dbPass = "1pwd2PWD3"
	DBName  string `xml:"db-name"`  // Oracle.bdName = "BDNAME"
	NLSLang string `xml:"nls-lang"` // Oracle.NLSLang = "RUSSIAN_RUSSIA.CL8MSWIN1251"
	DB      *sql.DB
}

//ConnectBD - .....
func (or *OracleConnect) ConnectBD() error {
	if Oracle.DB != nil {
		Oracle.DB.Close()
	}
	// defer db.Close() //  клозить будем на выходе из мэйна
	if len(or.NLSLang) > 0 {
		//os.Setenv("NLS_LANG", "RUSSIAN_RUSSIA.CL8MSWIN1251")
		os.Setenv("NLS_LANG", or.NLSLang)
	}
	orDB, err := sql.Open("godror", or.DBLogin+"/"+or.DBPass+"@"+or.DBHost)

	if err != nil {
		return errors.New("Проблемы с коннектом: " + or.DBLogin + "/****" + "@" + or.DBHost + "\n" + err.Error())
	}
	or.DB = orDB

	rows, err := or.DB.Query("select sysdate from dual")

	if err != nil {
		LogError := "(#51) Error running query: select sysdate from dual" + "\n" + err.Error()
		LogError += "(#52) " + or.DBHost + " ... " + or.DBName
		return errors.New(LogError)
	}
	defer rows.Close()

	var thedate string
	for rows.Next() {
		rows.Scan(&thedate)
	}
	return nil

}

func (or *OracleConnect) GetDBSettingFromXMLfile() error {

	// Open our Settings: oracle-settings.xml
	// пробуем открыть файл с настройками из каталога с программой
	ex, err := os.Executable()
	if err != nil {
		LogError := "(#settingsXML.getDBSettingFromXMLfile() провалилась попытка получить рабочую директорию (os.Executable()) :" + err.Error()
		return errors.New(LogError)
	}
	ExePath := filepath.Dir(ex)
	fileName := ExePath + "/oracle-settings.xml"
	xmlFile, err := os.Open(fileName)

	if err != nil {
		return errors.New("error os.Open(" + fileName + "):" + err.Error())
	}
	defer xmlFile.Close()

	byteXML, _ := ioutil.ReadAll(xmlFile)
	if err := xml.Unmarshal(byteXML, &or); err != nil {
		return errors.New("error xml.Unmarshal " + fileName + ":" + err.Error())
	}
	return nil
}
