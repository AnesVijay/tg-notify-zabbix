package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//function isProblem(params){
//    return params.Event_value == 1
//        && params.Event_update_status == 0
//    ;
//}

//function isUpdate(params){
//    return params.Event_value == 1
//        && params.Event_update_status == 1
//    ;
//}

func isResolution(eventValue string) bool{
	parsedInt, _ := strconv.ParseInt(eventValue, 10, 64)
    return parsedInt == 0
}

func getPriorityInMessageBySeverity(severity_level string) string {
        if (severity_level == "High"){
            return "High 🟧"
		} else if (severity_level == "Average"){
            return "Average 🟦"
		} else if (severity_level == "Disaster") {
            return "Disaster 🟥"
		} else {
            return "Minor ⬜"
		}
}

func main() {
	if len(os.Args) < 11 {
		log.Fatal("Usage: tg_notify <language: ru/en> <event_duration> <event_severity> <event_age> <event_value> <httpproxy> <message> <resolve_duration_min> <subject> <to> <token>") // <parsemode>
	}

	lang := os.Args[1]
	event_duration := os.Args[2]
	event_severity := os.Args[3]
	event_age := os.Args[4]
	event_value := os.Args[5]
	httpproxy := os.Args[6]
	message := os.Args[7]
	resolve_duration_min := os.Args[8]
	subject := os.Args[9]
	to := os.Args[10]
	token := os.Args[11]
	
//	parsemode := os.Args[11]
	if lang == "en" || lang == "ru" {
		SetLang(lang)
	}

	if token == "" {
		log.Fatal("Variable 'token' is not set")
	}

	chatID, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		log.Fatalf("Wrong Chat ID format: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Bot initialization error: %v", err)
	}

	if httpProxy := httpproxy; httpProxy != "" {
		proxyURL, err := url.Parse(httpProxy)
		if err != nil {
			log.Fatalf("Error parsing proxy URL: %v", err)
		}
		
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		
		// Change http client with set timeout
		bot.Client = &http.Client{
			Transport: transport,
			Timeout:   30,
		}
		log.Printf("Proxy is set: %s", httpProxy)
	}

	// Format a message into HTML
	text := fmt.Sprintf("<b>🚨 %s</b>\n\n<code>%s</code>", subject, message)

	// Optional (or for the future): it is possible to add inline-buttons
	/*
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗 Открыть в Zabbix", "https://zabbix.yourcompany.com/tr_events.php?triggerid=123"),
		),
	)
	*/

	// MODIFYING MESSAGE
	modified_text := strings.Replace(text,CurrentLang[Severity]+":",CurrentLang[Severity] + ": " + getPriorityInMessageBySeverity(event_severity),1)

	if (strings.Contains(modified_text, strings.ToUpper(CurrentLang[Problem]))){
        modified_text = "🚨 " + modified_text;
    } else if (strings.Contains(modified_text,strings.ToUpper(CurrentLang[Problem]))) {
        modified_text = "✅ " + modified_text
    } else if (strings.Contains(modified_text,strings.ToUpper(CurrentLang[Problem]))) {
        modified_text = "🔄 " + modified_text;
    }

    // РЕШЕНИЕ - при решении события
    if (isResolution(event_value)) {
        if (event_age < resolve_duration_minutes) {
            log.Fatalf("Event duration(%s) is less than DELAY (%sm)", event_duration, resolve_duration_min)
        }
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"


	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}
}
