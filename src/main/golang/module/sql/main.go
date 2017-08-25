package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"os"

	"path/filepath"
	"runtime"

	"path"

	sql2 "main/golang/util/sql"
	u "main/golang/util/string"
)


func main() {

	sql, _ := ioutil.ReadFile(getCurrentFilePath() + "sql")
	var Sql = string(sql)
	//fmt.Println(Sql)
	var buffer bytes.Buffer

	tableName := regexp.MustCompile(`create table (.*\.)?(.*)`).FindStringSubmatch(Sql)
	commentReg := regexp.MustCompile("comment\\s+'(.*)'\\s*")
	if len(tableName) < 2 {
		log.Println(tableName)
		return
	}

	var class_name = u.ToUpperCamel(tableName[2])

	buffer.WriteString("public class " + class_name + "{ \n")
	//log.Println(class_name)
	content := Sql[strings.Index(Sql, "(")+1: strings.LastIndex(Sql, ")")]
	//log.Println(content)
	lines := strings.Split(content, "\n")
	var field, comments []string
	var strType, comment string
	for _, line := range lines {
		//log.Println("i " + strconv.Itoa(i) + "  " + line)

		field = strings.Fields(line)
		comments = commentReg.FindStringSubmatch(line)
		if len(comments) > 0 {
			comment = commentReg.FindStringSubmatch(line)[1]
			buffer.WriteString("    //  " + comment + "\n")
		}

		if len(field) > 2 {
			strType = sql2.GetJavaTypeByMySql(field[1])
			//log.Println(field)
			buffer.WriteString("    private " + strType + " " + u.ToLowerCamel(field[0]) + ";\n")
		}
	}

	buffer.WriteString("}\n")
	file_name := class_name + "df"
	log.Println(file_name)
	var java_file string = path.Join("main/golang/module/sql/", class_name+"java")
	file, err := os.Create(java_file)

	if err != nil {
		log.Fatal(err)
	} else {
		file.WriteString(buffer.String())
	}

}

func getCurrentFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir1, _ := filepath.Split(filename)
	return dir1
}


