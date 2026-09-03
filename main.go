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

//function isResolution(params){
//    return params.Event_value == 0;
//}

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
		log.Fatal("Использование: tg_notifications <event_duration> <event_severity> <event_update_status> <event_value> <httpproxy> <message> <resolve_duration_min> <subject> <to> <token>") // <parsemode>
	}

	event_duration := os.Args[1]
	event_severity := os.Args[2]
	event_update_status := os.Args[3]
	event_value := os.Args[4]
	httpproxy := os.Args[5]
	message := os.Args[6]
	resolve_duration_min := os.Args[7]
	subject := os.Args[8]
	to := os.Args[9]
	token := os.Args[10]
//	parsemode := os.Args[11]

	if token == "" {
		log.Fatal("Переменная token не установлена")
	}

	// Парсим Chat ID (он может быть отрицательным для групп/каналов)
	chatID, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		log.Fatalf("Неверный формат Chat ID: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}

	if httpProxy := httpproxy; httpProxy != "" {
		proxyURL, err := url.Parse(httpProxy)
		if err != nil {
			log.Fatalf("Ошибка парсинга URL прокси: %v", err)
		}
		
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		
		// Подменяем клиент, сохраняя стандартный таймаут библиотеки
		bot.Client = &http.Client{
			Transport: transport,
			Timeout:   30,
		}
		log.Printf("Настроен прокси: %s", httpProxy)
	}

	// Формируем красивое сообщение с HTML-разметкой
	text := fmt.Sprintf("<b>🚨 %s</b>\n\n<code>%s</code>", subject, message)

	// Опционально: можно добавить inline-кнопки прямо из Zabbix-скрипта!
	/*
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗 Открыть в Zabbix", "https://zabbix.yourcompany.com/tr_events.php?triggerid=123"),
		),
	)
	*/

	// ОБРАБОТКА СООБЩЕНИЯ
	modified_text := strings.Replace(text,"Серьёзность:", "Серьёзность: " + getPriorityInMessageBySeverity(event_severity),1)

	if (strings.Contains(modified_text, "ПРОБЛЕМА")){
        modified_text = "🚨 " + modified_text;
    } else if (strings.Contains(modified_text,"РЕШЕНО")) {
        modified_text = "✅ " + modified_text
    } else if (strings.Contains(modified_text,"ОБНОВЛЕНА")) {
        modified_text = "🔄 " + modified_text;
    }

    // РЕШЕНИЕ - при решении события
    if (isResolution(params)) {
        if (parseZabbixTimeToSeconds(params.Event_duration) >= parseZabbixTimeToSeconds(params.Resolve_duration_minutes)) {
            Telegram.sendMessage();
        } else {
            log.Fatalf("Event duration(%s) is less than DELAY (%sm)", event_duration, resolve_duration_min);
        }
    }

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"


	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}
}
