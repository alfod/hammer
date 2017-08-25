package sql

import (
	"strings"
	"log"
)

var mysql_java_type_map = map[string]string{
	"VARCHAR":   "String",
	"CHAR":      "String",
	"BLOB":      "byte",
	"TEXT":      "String",
	"INTEGER":   "Long",
	"TINYINT":   "Integer",
	"SMALLINT":  "Integer",
	"MEDIUMINT": "Integer",
	"BIT":       "Boolean",
	"BIGINT":    "BigInteger",
	"FLOAT":     "Float",
	"DOUBLE":    "Double",
	"DECIMAL":   "BigDecimal",
	"BOOLEAN":   "Integer",
	"ID":        "Long",
	"DATE":      "Date",
	"TIME":      "Time",
	"DATETIME":  "Date",
	"TIMESTAMP": "Date",
	"YEAR":      "Date",
	"INT":       "Integer",
}

func GetJavaTypeByMySql(mysql_type string) string{
	if strings.Contains(mysql_type, "(") {
		mysql_type = mysql_type[:strings.Index(mysql_type, "(")]
	}
	var str = mysql_java_type_map[strings.TrimSpace(strings.ToUpper(mysql_type))]
	if str == "" {
		//fmt.Println(mysql_type + " is not support ")
		log.Fatal(mysql_type + " is not support ")
		//return str, errors.New("not support type " + mysql_type)
	}
	return str
}
