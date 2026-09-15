package main

type WordType int

const (
	Severity WordType = iota
	Problem
	Updated
	Resolved
)

var ruLang = map[WordType]string {
	Severity : "Серьёзность",
	Problem : "Проблема",
	Updated : "Обновление",
	Resolved : "Решено",
}

var enLang = map[WordType]string {
	Severity : "Severity",
	Problem : "Problem",
	Updated : "Update",
	Resolved : "Resolved",
}

var CurrentLang map[WordType]string

func SetLang(langCode string) {
	if (langCode == "en") {
		CurrentLang = enLang
	} else {
		CurrentLang = ruLang
	}
}
