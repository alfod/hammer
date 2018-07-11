package java

var java_package_map = map[string]string{
	"Date": "import java.util.Date;",
}

func GetJavaPackageByType(java_type string) string {
	var str = java_package_map[java_type]
	return str
}
