package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"os"

	sql2 "main/golang/util/sql"

	"path/filepath"
	"runtime"

	strings2 "main/golang/util/string"
	"strconv"
)


func main() {

	sqlBytes, _ := ioutil.ReadFile(GetCurrentFilePath() + "sql")
	sqlBytes = bytes.Replace(sqlBytes, []byte("\r"), []byte(""), -1)
	var sqlStr = string(sqlBytes)
	//log.Println(sqlStr)
	//var sqlTest []string = regexp.MustCompile(`(\s*create\s+table.+\s*\n?\((.*\n)+)`).FindStringSubmatch(sqlStr)
	//log.Println(sqlTest)
	//var sqlTest []string = regexp.MustCompile(`(.*\n)+`).FindStringSubmatch(sqlStr)
	//log.Println(sqlTest)
	var sqls []string = regexp.MustCompile(`(\s*create{1}\s+table.+\s*\n?\((.*\n)+?\)\n)+`).FindStringSubmatch(sqlStr)
	//log.Println(sqls)
	if len(sqls) > 0 {
		for i, sql := range sqls {
			log.Println("i: "+strconv.Itoa(i)+"\n"+sql)
			dealSingleCreateSql(sql)
		}
	} else {
		log.Fatal("split file failed")
	}

}

func dealSingleCreateSql(sql string) {
	var buffer bytes.Buffer
	tableName := regexp.MustCompile(`create\s+table\s+(\S+\.)?(\S+)`).FindStringSubmatch(sql)
	commentReg := regexp.MustCompile(`comment\s+'(\S*)'\s*`)
	if len(tableName) < 2 {
		log.Fatal(tableName)
	}
	var class_name = strings2.ToUpperCamel(tableName[2])

	buffer.WriteString("public class " + class_name + "{ \n")
	content := sql[strings.Index(sql, "(")+1: strings.LastIndex(sql, ")")]

	var lines []string
	if strings.Contains(content, ",\r\n") {
		lines = strings.Split(content, ",\r\n")
	} else {
		lines = strings.Split(content, ",\n")
	}

	var field, comments []string
	var strType, comment string
	for _, line := range lines {
		field = strings.Fields(line)
		comments = commentReg.FindStringSubmatch(line)
		if len(comments) > 0 {
			comment = comments[1]
			buffer.WriteString("    //*  " + comment + "*/ \n")
		}
		if len(field) > 2 {
			strType = sql2.GetJavaTypeByMySql(field[1])
			buffer.WriteString("    private " + strType + " " + strings2.ToLowerCamel(field[0]) + ";\n")
		}
	}

	buffer.WriteString("}\n")
	var java_file string = "main/golang/module/sql/" + class_name + ".java"
	file, err := os.Create(java_file)
	if err != nil {
		log.Fatal(err)
	} else {
		file.WriteString(buffer.String())
	}
	file.Close()
}

func GetCurrentFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir1, _ := filepath.Split(filename)
	return dir1
}




