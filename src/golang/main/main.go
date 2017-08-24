package main

import (
	"fmt"
	"regexp"
	"strings"
	u "golang/util/string"
)

var sql = `
create table institution.ins_sign_info
(
	id int
	auto_increment primary key,
	institution_id int not null comment '机构ID',
	manage_id int null comment '客户经理ID',
	standard_amount double(10,2) null comment '标准最低应付',
	discount_amount double(10,2) null comment '折后最低应付',
	discount double(10,2) null comment '签约折扣',
	sign_amount double(11,2) null comment '合同签约价',
	status int(2) null comment '合同/协议状态  :0:未签合同  1:已签合同  2：已同意协议',
	school_license varchar(100) null comment '办学许可',
	business_license varchar(100) null comment '营业执照',
	legal_person_id_card varchar(100) null,
	apply_time datetime null comment '申请时间',
	sign_code varchar(100) null comment '交易编号',
	agreement_id int null comment '协议',
	agreement_params text null comment '协议参数列表{ key1: value1 , key2 :  vlaue2}',
	sign_start_date date null comment '签约开始日期',
	sign_end_date date null comment '签约结束日期',
	manager_phone varchar(20) null comment '管理员手机',
	manager_name varchar(50) null comment '管理员姓名',
	manager_name_pinyin varchar(255) null comment '管理员姓名拼音',
	ali_sycn_params text null comment '支付宝同步推送参数',
	ali_asyn_params text null comment '支付宝异步步推送参数',
	payed_amount double(10,2) null,
	pay_status int(2) null comment '支付状态 0 未支付 1 线上支付 2 线下支付',
	lastupdate bigint(13) unsigned zerofill null,
	is_intent tinyint(1) default '0' null comment '是否交意向金 0否 1是',
	intent_money double(20,2) default '0.00' null comment '意向金额',
	disclaimer_license varchar(100) null,
	pay_model varchar(50) null comment '付款方式（转账1  刷卡2）',
	deleted tinyint default '0' not null comment '0 没有删除， 1 已经删除',
	update_time datetime default '1970-01-01 00:00:00' not null comment '修改日期',
	create_time datetime default '1970-01-01 00:00:00' not null comment '创建日期',
	pay_comment varchar(150) null comment '汇款信息备注'
)
`

func main() {
	tableName := regexp.MustCompile(`create table [\w\d.]*\.?(.*)`).FindStringSubmatch(sql)
	commentReg := regexp.MustCompile("comment\\s+'(.*)'\\s*")
	fmt.Println("public class " + u.UnderLineToCamel(tableName[1]) + "{")
	firstBracket := strings.Index(sql, "(")
	lastBracket := strings.LastIndex(sql, ")")
	content := sql[firstBracket+1: lastBracket]
	//fmt.Println(content)
	lines := strings.Split(content, ",\n")
	var field []string
	var strType string
	var comment string
	var comments []string
	for _, line := range lines {
		field = strings.Fields(line)
		comments = commentReg.FindStringSubmatch(line)
		if len(comments) > 0 {
			comment = commentReg.FindStringSubmatch(line)[1]
			fmt.Println("//  " + comment)
		}

		if len(field) > 2 {
			if strings.Contains(field[1], "int") {
				strType = "Integer"
			} else if strings.Contains(field[1], "char") ||
				strings.Contains(field[1], "text") {
				strType = "String"
			} else if strings.Contains(field[1], "double") {
				strType = "Double"
			} else if strings.Contains(field[1], "date") {
				strType = "Date"
			} else {
				fmt.Println()
				fmt.Println(line)
				panic("not support field type " + field[1])
				fmt.Println()

			}
			fmt.Println("private " + strType + " " + u.ToCamel(field[0], false) + ";")
		}
	}

	fmt.Println("}")
}

type Field struct {
	name  string
	genre string
}
