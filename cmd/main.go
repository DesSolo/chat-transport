package main

import (
	"chat-transport/internal/config"
	"chat-transport/internal/daemon"
	"chat-transport/internal/entities"
	"chat-transport/internal/transport/telegram"
	"chat-transport/internal/transport/vk"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"
)

var configFile string

var version string

// GetTransports ...
func GetTransports(chats map[string]*config.Chat) ([]entities.Transport, error) {
	var transports []entities.Transport

	for _, chat := range chats {

		switch chat.Type {
		case "telegram":
			tpl, err := template.New("message").Parse(chat.Template)
			if err != nil {
				log.Fatal(err)
			}

			tg := telegram.NewTelegram(chat.Name, chat.Token, chat.ChatID, chat.IgnoreAccounts, tpl)
			transports = append(transports, tg)

		case "vk":
			chatID, err := strconv.Atoi(chat.ChatID)
			if err != nil {
				return nil, errors.New("chat_id must be integer")
			}

			var ignoreAccounts []int

			for _, sid := range chat.IgnoreAccounts {
				iid, err := strconv.Atoi(sid)
				if err != nil {
					return nil, errors.New("ignore account id must be integer")
				} 

				ignoreAccounts = append(ignoreAccounts, iid)
			}

			vk := vk.NewClient(chat.Name, chat.Token, chat.Cookies["remixsid"], chatID, ignoreAccounts)
			transports = append(transports, vk)

		default:
			return nil, fmt.Errorf("transport \"%s\" not supported", chat.Type)
		}
	}

	for _, tr := range transports {
		if err := tr.Validate(); err != nil {
			return nil, err
		}
	}

	return transports, nil
}

func main() {
	flag.StringVar(&configFile, "c", "config.toml", "config file path")
	showVer := flag.Bool("v", false, "show current version")
	flag.Parse()

	if *showVer {
		fmt.Printf("chat transport version: %s\n", version)
		os.Exit(0)
	}

	conf, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("fault read config file err: \"%s\"", err)
	}

	srcTransports, err := GetTransports(conf.Src)
	if err != nil {
		log.Fatalf("fault load srt transports err: %s", err)
	}

	log.Printf("success loaded %d src chat(s)", len(srcTransports))

	dstTransports, err := GetTransports(conf.Dst)
	if err != nil {
		log.Fatalf("fault load dst transports err: %s", err)
	}

	log.Printf("success loaded %d dst chat(s)", len(srcTransports))

	updateInterval := conf.Interval * time.Second

	daemon := daemon.NewDaemon(srcTransports, dstTransports, updateInterval)

	log.Printf("consume updates every %s", updateInterval)

	if err := daemon.Run(); err != nil {
		log.Fatal(err)
	}
}
