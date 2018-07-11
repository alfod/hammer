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
	sqlKey := `\s*KEY|key\s+.*\n`;
	sqlKeyPattern, _ := regexp.Compile(sqlKey)
	sqlStr = sqlKeyPattern.ReplaceAllString(sqlStr, "")
	//log.Println(sqlStr)
	var sqls []string = regexp.MustCompile(`(\s*create\s+table\s+\w+\s*\n?\(\s*\n?(\s*[a-zA-Z\']+.*\n?)+\s*\n*\s*\)\s*\n*)+`).FindAllString(sqlStr, -1)
	if len(sqls) > 0 {
		for i, sql := range sqls {
			log.Println("i: " + strconv.Itoa(i) + "\n")
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
	content := sql[strings.Index(sql, "(")+1 : strings.LastIndex(sql, ")")]

	var lines []string
	if strings.Contains(content, ",\r\n") {
		lines = strings.Split(content, ",\r\n")
	} else {
		lines = strings.Split(content, ",\n")
	}

	var lineWords, comments []string
	var comment string
	var fields = make([]string, len(lines))
	var types = make([]string, len(lines))
	for index, line := range lines {
		lineWords = strings.Fields(line)
		if len(lineWords) < 2 {
			continue
		}
		comments = commentReg.FindStringSubmatch(line)
		if len(comments) > 0 {
			comment = comments[1]
			buffer.WriteString("     /** \n")
			buffer.WriteString("      *  " + comment + "\n")
			buffer.WriteString("      */ \n")
		}
		fields[index] = strings2.ToLowerCamel(lineWords[0])
		types[index] = sql2.GetJavaTypeByMySql(lineWords[1])
		if len(fields[index]) > 2 {
			buffer.WriteString("     private " + types[index] + " " + fields[index] + ";\n")
		}
	}
	var upperField string
	var byteString []byte
	buffer.WriteString(" \n")
	for index, field := range fields {
		if len(field) < 1 {
			continue
		}
		byteString = []byte(field)
		byteString[0] = field[0] - 32
		upperField = string(byteString)
		//setter
		buffer.WriteString("     public  void set" + upperField + " (" + types[index] + " " + field + "){\n")
		buffer.WriteString("            this." + field + " = " + field + ";\n")
		buffer.WriteString("     } \n")
		//getter
		buffer.WriteString("     public  void get" + upperField + " (){\n")
		buffer.WriteString("            return this." + field + ";\n ")
		buffer.WriteString("    } \n")
	}

	buffer.WriteString("\n")
	buffer.WriteString(`     @Override
     public String toString() {
		`)
	buffer.WriteString("   return \"" + class_name + "{\" + \n")
	for index, field := range fields {
		if len(field) < 1 {
			continue
		}
		buffer.WriteString("                     \"")
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(field + " = \" + " + field + " + ")
		if index >= len(fields)-2 {
			buffer.WriteString(" \"}\";")
		}
		buffer.WriteString(" \n ")
	}
	buffer.WriteString("\n   	  }")
	buffer.WriteString("\n }")

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
