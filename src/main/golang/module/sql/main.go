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

	"../../util/java"
	strings2 "../../util/string"
)

func main() {

	sqlBytes, _ := ioutil.ReadFile(GetCurrentFilePath() + "sql")
	sqlBytes = bytes.Replace(sqlBytes, []byte("\r"), []byte(""), -1)
	sqlBytes = bytes.Replace(sqlBytes, []byte("`"), []byte(""), -1)
	var sqlStr = strings.ToLower(string(sqlBytes))
	sqlPrimarykey := `\s*primary key.*\n`;
	sqlPrimaryKeyPattern, _ := regexp.Compile(sqlPrimarykey)
	sqlStr = sqlPrimaryKeyPattern.ReplaceAllString(sqlStr, "")

	sqlKey := `\s*\bkey\b.*\n`;
	sqlKeyPattern, _ := regexp.Compile(sqlKey)
	sqlStr = sqlKeyPattern.ReplaceAllString(sqlStr, "")

	sqlBracket := `\(\d+\)`;
	sqlBracketPattern, _ := regexp.Compile(sqlBracket)
	sqlStr = sqlBracketPattern.ReplaceAllString(sqlStr, "")
	//log.Println(sqlStr)
	var sqls []string = regexp.MustCompile(`(\s*create\s+table\s+\w+\s*\n?\(\s*\n?(\s*[a-zA-Z\'\,]+.*\n?)+\s*\n*\s*\)\s*\n*)+`).FindAllString(sqlStr, -1)
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
	buffer.WriteString("\n\n")
	content := sql[strings.Index(sql, "(")+1 : strings.LastIndex(sql, ")")]
	var lines []string
	if strings.Contains(content, ",\r\n") {
		lines = strings.Split(content, ",\r\n")
	} else {
		lines = strings.Split(content, ",\n")
	}
	tableName := regexp.MustCompile(`\s*create\s+table\s+(\w+\.)?(\w+)`).FindStringSubmatch(sql)
	if len(tableName) < 2 {
		log.Fatal(tableName)
	}
	commentReg := regexp.MustCompile(`comment\s+'(.*)'\s*`)
	var lineWords, originComments []string
	var comment, packagePath string
	var fields = make([]string, len(lines))
	var comments = make([]string, len(lines))
	var types = make([]string, len(lines))
	var packagePathMap = map[string]string{}
	for index, line := range lines {
		lineWords = strings.Fields(line)
		if len(lineWords) < 2 {
			continue
		}
		fields[index] = strings2.ToLowerCamel(lineWords[0])
		types[index] = sql2.GetJavaTypeByMySql(lineWords[1])
		originComments = commentReg.FindStringSubmatch(line)
		if len(originComments) > 0 {
			comments[index] = originComments[1]
		}

		if (len(types[index]) > 1) {
			packagePath = java.GetJavaPackageByType(types[index])
			var result = packagePathMap[types[index]]
			if (packagePath != "" && result == "") {
				packagePathMap[types[index]] = packagePath
				buffer.WriteString(packagePath)

			}
		}
	}
	buffer.WriteString("\n\n")
	var class_name = strings2.ToUpperCamel(tableName[2])

	buffer.WriteString("public class " + class_name + "{ \n")

	for index, field := range fields {
		if len(field) < 1 {
			continue
		}
		comment = comments[index]
		if len(comment) > 0 {
			buffer.WriteString("     /** \n")
			buffer.WriteString("      *  " + comment + "\n")
			buffer.WriteString("      */ \n")
		}
		if len(field) > 1 {
			buffer.WriteString("     private " + types[index] + " " + field + ";\n")
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
		buffer.WriteString("     public  void set" + upperField + " (" + types[index] + " " + field + ") {\n")
		buffer.WriteString("            this." + field + " = " + field + ";\n")
		buffer.WriteString("     } \n")
		//getter
		buffer.WriteString("     public " + types[index] + " get" + upperField + " () {\n")
		buffer.WriteString("            return this." + field + ";\n ")
		buffer.WriteString("    } \n")
	}

	buffer.WriteString("\n")
	//toString
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
		buffer.WriteString(" \n ")
	}
	buffer.WriteString(" 				\"}\";")
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
