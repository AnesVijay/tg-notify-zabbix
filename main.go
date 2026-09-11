package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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
	if len(os.Args) < 13 {
		log.Fatal("Usage: tg_notify <language: ru/en> <event_duration> <event_severity> <event_timestamp> <event_value> <httpproxy> <message> <resolve_duration_sec> <subject> <to> <token> <proxy_timeout>") // <parsemode>
	}

	lang := os.Args[1]
	eventDuration := os.Args[2]
	eventSeverity := os.Args[3]
	eventTimestampInp := os.Args[4]
	eventValue := os.Args[5]
	httpProxy := os.Args[6]
	message := os.Args[7]
	resolveDurationSecInp := os.Args[8]
	subject := os.Args[9]
	to := os.Args[10]
	token := os.Args[11]
	proxyTimeout := os.Args[12]

	var botClient *http.Client
	var bot *tgbotapi.BotAPI
	
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


	if httpProxy != "" {
		proxyStr := httpProxy
		
		// "http://" as a default scheme 
		// This prevents url.Parse from treating the hostname as the scheme.
		if !strings.HasPrefix(proxyStr, "http://") && 
		   !strings.HasPrefix(proxyStr, "https://") && 
		   !strings.HasPrefix(proxyStr, "socks5://") {
			proxyStr = "http://" + proxyStr
		}

		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Fatalf("Error parsing proxy URL: %v", err)
		}

		timeout, err := strconv.ParseInt(proxyTimeout, 10, 64)
		if err != nil {
			log.Fatalf("Error parsing proxy_timeout: %v", err)
		}
		
		// connection pooling, TLS config, dialers instead of a bare struct.
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.Proxy = http.ProxyURL(proxyURL)
		
		botClient = &http.Client{
			Transport: transport,
			Timeout:   time.Duration(timeout) * time.Second,
		}
		log.Printf("Proxy is set: %s", proxyStr)
		bot, err = tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, botClient)
		if err != nil {
			log.Fatalf("Bot initialization error: %v", err)
		}
	} else {
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Fatalf("Bot initialization error: %v", err)
		}
	}

	eventTimestamp, err := strconv.ParseInt(eventTimestampInp, 11, 64)
	if err != nil {
		log.Fatalf("Error parsing event_timestamp (after severity): %v", err)
	}
	resolveDurationSec, err := strconv.ParseInt(resolveDurationSecInp, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing resolve_duration_sec (after message): %v", err)
	}


	// Format a message into HTML
	text := fmt.Sprintf("<b>%s</b>\n\n<code>%s</code>", subject, message)

	// Optional (or for the future): it is possible to add inline-buttons
	/*
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗 Открыть в Zabbix", "https://zabbix.yourcompany.com/tr_events.php?triggerid=123"),
		),
	)
	*/

	// MODIFYING MESSAGE
	// adding emoji to severity
	modifiedText := strings.Replace(text,CurrentLang[Severity]+":",CurrentLang[Severity] + ": " + getPriorityInMessageBySeverity(eventSeverity),1)
	// adding emoji to illustarte event state
	if (strings.Contains(modifiedText, strings.ToUpper(CurrentLang[Problem]))){
        modifiedText = "🚨 " + modifiedText;
    } else if (strings.Contains(modifiedText,strings.ToUpper(CurrentLang[Resolved]))) {
        modifiedText = "✅ " + modifiedText
    } else if (strings.Contains(modifiedText,strings.ToUpper(CurrentLang[Updated]))) {
        modifiedText = "🔄 " + modifiedText;
    }

    // in case of RESOLVED check event duration to avoid flapping alerts
    if (isResolution(eventValue)) {
		eventTime := time.Unix(eventTimestamp, 0)
		durationSeconds := int64(time.Since(eventTime).Seconds())
        if (durationSeconds < resolveDurationSec) {
            log.Fatalf("Event duration(%s) is less than DELAY (%d seconds)", eventDuration, resolveDurationSec)
        }
	}

	msg := tgbotapi.NewMessage(chatID, modifiedText)
	msg.ParseMode = "HTML"


	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Error sending a message: %v", err)
	}
}
