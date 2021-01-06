package orc

import (
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/alex-mj-pradius/fox-log"
	_ "github.com/godror/godror"
)

//Log to all function
var Log log.Log

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
	DB      *sql.DB
}

//ConnectBD - .....
func (or *OracleConnect) ConnectBD() {
	if Oracle.DB != nil {
		Oracle.DB.Close()
	}

	Log.Info("\n******** Поднять коннекты к базам ********")
	os.Setenv("NLS_LANG", "RUSSIAN_RUSSIA.CL8MSWIN1251")
	orDB, err := sql.Open("godror", or.DBLogin+"/"+or.DBPass+"@"+or.DBHost)

	if err != nil {
		Log.Error("Проблемы с коннектом: " + or.DBLogin + "/****" + "@" + or.DBHost + "\n" + err.Error())
		return
	}
	or.DB = orDB

	rows, err := or.DB.Query("select sysdate from dual")

	if err != nil {
		Log.Error("(#87) !!! Error running query: select sysdate from dual" + "\n" + err.Error())
		Log.Error("(#88) " + or.DBHost + " ... " + or.DBName)
		return
	}
	defer rows.Close()

	var thedate string
	for rows.Next() {
		rows.Scan(&thedate)
	}
	Log.Info("*** test to: " + or.DBHost + "  ***")
	Log.Info("*** date is: " + thedate + " ***")

	Log.Debug("******************************************")
	// defer db.Close() //  клозить будем на выходе из мэйна
}

func (or *OracleConnect) GetDBSettingFromXMLfile() {

	// Open our Settings: oracle-settings.xml
	// пробуем открыть файл с настройками из каталога с программой
	ex, err := os.Executable()
	if err != nil {
		Log.Error("(#settingsXML.getDBSettingFromXMLfile() провалилась попытка получить рабочую директорию (os.Executable()) :" + err.Error())
		panic(err)
	}
	ExePath := filepath.Dir(ex)

	settingFileName := "oracle-settings.xml"
	xmlFile, err := os.Open(ExePath + "/" + settingFileName)

	if err != nil {
		Log.Debug(err.Error())
		// не получилось открываем константный, отладка
		//debugFileName := "/home/parallels/go/src/" + settingFileName
		debugFileName := "c:\\Go\\src\\" + settingFileName
		Log.Debug("... " + settingFileName + " не найден в папке с исполняемым файлом, открываем " + debugFileName)
		xmlFile, err = os.Open(debugFileName)
		if err != nil {
			Log.Error("error os.Open(" + debugFileName + "):" + err.Error())
		}
	}
	defer xmlFile.Close()

	byteXML, _ := ioutil.ReadAll(xmlFile)
	if err := xml.Unmarshal(byteXML, &or); err != nil {
		Log.Error("error xml.Unmarshal " + settingFileName + ":" + err.Error())
		panic(err)
	}
}
