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

func getPriorityInMessageBySeverity(severityLevel string) string {
    switch severityLevel {
	case "High":
		return "High 🟧"
	case "Average":
		return "Average 🟦"
	case "Disaster":
		return "Disaster 🟥"
	default:
		return "Minor ⬜"
	}
}

func main() {
	if len(os.Args) < 11 {
		log.Fatal("Usage: tg_notify <language: ru/en> <event_duration> <event_severity> <event_age> <event_value> <httpproxy> <message> <resolve_duration_sec> <subject> <to> <token>") // <parsemode>
	}

	lang := os.Args[1]
	eventDuration := os.Args[2]
	eventSeverity := os.Args[3]
	eventAge := os.Args[4]
	eventValue := os.Args[5]
	httpProxy := os.Args[6]
	message := os.Args[7]
	resolveDurationSec := os.Args[8]
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

	if httpProxy := httpProxy; httpProxy != "" {
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
	// add emoji to severity
	modifiedText := strings.Replace(text,CurrentLang[Severity]+":",CurrentLang[Severity] + ": " + getPriorityInMessageBySeverity(eventSeverity),1)
	// add emoji to illustarte event state
	if (strings.Contains(modifiedText, strings.ToUpper(CurrentLang[Problem]))){
        modifiedText = "🚨 " + modifiedText;
    } else if (strings.Contains(modifiedText,strings.ToUpper(CurrentLang[Resolved]))) {
        modifiedText = "✅ " + modifiedText
    } else if (strings.Contains(modifiedText,strings.ToUpper(CurrentLang[Updated]))) {
        modifiedText = "🔄 " + modifiedText;
    }

    // in case of RESOLVED
    if (isResolution(eventValue)) {
        if (eventAge < resolveDurationSec) {
            log.Fatalf("Event duration(%s) is less than DELAY (%sm)", eventDuration, resolveDurationSec)
        }
	}

	msg := tgbotapi.NewMessage(chatID, modifiedText)
	msg.ParseMode = "HTML"


	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}
}
