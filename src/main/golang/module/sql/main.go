package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"os"
    "strconv"
	sql2 "../../util/sql"

	"path/filepath"
	"runtime"

	strings2 "../../util/string"
)


func main() {

	sqlBytes, _ := ioutil.ReadFile(GetCurrentFilePath() + "sql")
	sqlBytes = bytes.Replace(sqlBytes, []byte("\r"), []byte(""), -1)
	sqlBytes = bytes.Replace(sqlBytes, []byte("`"), []byte(""), -1)
	var sqlStr = strings.ToLower(string(sqlBytes))
	//log.Println(sqlStr)
	var sqls []string = regexp.MustCompile(`(\s*create\s+table\s+\w+\s*\n?\(\s*\n?(\s*[a-zA-Z\']+.*\n?)+\s*\n*\s*\)\s*\n*)+`).FindAllString(sqlStr, -1)
	if len(sqls) > 0 {
		for i,sql := range sqls {
			log.Println("i: "+strconv.Itoa(i)+"\n")
			dealSingleCreateSql(sql)
		}
	} else {
		log.Fatal("split file failed")
	}

}

func dealSingleCreateSql(sql string) {
	var buffer bytes.Buffer
	tableName := regexp.MustCompile(`create\s+table\s+(\S+\.)?(\S+)`).FindStringSubmatch(sql)
	commentReg := regexp.MustCompile(`comment\s+'(.*)'\s*`)
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
			buffer.WriteString("    /**  " + comment + "*/ \n")
		}
		if len(field) > 2 {
			strType = sql2.GetJavaTypeByMySql(field[1])
			buffer.WriteString("    private " + strType + " " + strings2.ToLowerCamel(field[0]) + ";\n")
		}
	}

	buffer.WriteString("}\n")
	var java_file string = "./" + class_name + ".java"
	file, err := os.Create(java_file)
	file_content := buffer.String()
	if err != nil {
		log.Fatal(err)
	} else {
		file.WriteString(file_content)
	}
	file.Close()
}

func GetCurrentFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir1, _ := filepath.Split(filename)
	return dir1
}




