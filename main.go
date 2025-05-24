package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

//go:embed banner.txt
var appBanner string
var appVersion string = "1.0.2"

type stringSlice []string

func (s *stringSlice) String() string     { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error { *s = append(*s, v); return nil }

var verbose bool

func logInfo(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func logError(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func logWarn(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

func logDebug(format string, args ...interface{}) {
	if verbose {
		log.Printf("[DEBUG] "+format, args...)
	}
}

func logFatal(format string, args ...interface{}) {
	log.Fatalf("[FATAL] "+format, args...)
}

func main() {
	var appID int
	var appHash, sessionsDir, phonesFile string
	var channels stringSlice
	var help, version bool

	flag.IntVar(&appID, "app-id", 0, "Telegram App ID")
	flag.StringVar(&appHash, "app-hash", "", "Telegram App Hash")
	flag.StringVar(&sessionsDir, "sessions-dir", "sessions", "Sessions directory")
	flag.StringVar(&phonesFile, "phones-file", "phones.txt", "Phone numbers file")
	flag.Var(&channels, "channel", "Channel to join (repeatable)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&version, "version", false, "Show version")

	flag.Usage = func() {
		app := filepath.Base(os.Args[0])
		fmt.Printf("USAGE: %s [OPTIONS]\n\nREQUIRED:\n  -app-id <int>      Telegram App ID\n  -app-hash <string> Telegram App Hash\n\nOPTIONS:\n  -sessions-dir <dir>   Session files directory (default: sessions)\n  -phones-file <file>   Phone numbers file (default: phones.txt)\n  -channel <channel>    Channel to join (repeatable)\n  -verbose              Detailed logging\n  -help                 Show this help\n  -version              Show version\n\nEXAMPLES:\n  %s -app-id 12345 -app-hash \"abc123\"\n  %s -app-id 12345 -app-hash \"abc123\" -channel @mychannel -verbose\n", app, app, app)
	}
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Ltime)
	if !verbose {
		log.SetFlags(0)
	}

	if help {
		flag.Usage()
		return
	}
	if version {
		fmt.Println(appBanner)
		fmt.Printf("sgen v.%v\nSemi-automated Telegram session creation tool (via Gogram)\n", appVersion)
		return
	}

	if appID == 0 || appHash == "" {
		fmt.Fprintf(os.Stderr, "Error: -app-id and -app-hash are required\nUse '%s -help' for usage.\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	f, err := os.Open(phonesFile)
	if err != nil {
		logFatal("Failed to open phones file: %v", err)
	}
	defer f.Close()

	var phones []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "//") {
			phones = append(phones, strings.TrimPrefix(line, "+"))
		}
	}
	if len(phones) == 0 {
		logFatal("No valid phones found")
	}

	os.MkdirAll(sessionsDir, 0755)
	logInfo("Starting with %d phones and %d channels", len(phones), len(channels))

	success := 0
	for i, phone := range phones {
		logInfo("Processing phone %d/%d: %s", i+1, len(phones), phone)

		client, err := telegram.NewClient(telegram.ClientConfig{
			AppID:    int32(appID),
			AppHash:  appHash,
			LogLevel: telegram.LogDisable,
			Session:  filepath.Join(sessionsDir, phone+".dat"),
		})
		if err != nil {
			logError("Client creation failed: %v", err)
			continue
		}

		_, err = client.Login(phone, &telegram.LoginOptions{
			CodeCallback: func() (string, error) {
				fmt.Printf("Enter code for %s: ", phone)
				var code string
				fmt.Scanln(&code)
				return strings.TrimSpace(code), nil
			},
			PasswordCallback: func() (string, error) {
				fmt.Printf("Enter 2FA for %s: ", phone)
				var pwd string
				fmt.Scanln(&pwd)
				return strings.TrimSpace(pwd), nil
			},
		})
		if err != nil {
			logError("Auth failed: %v", err)
			continue
		}

		logInfo("Auth successful")
		success++

		for _, ch := range channels {
			if _, err := client.JoinChannel(ch); err != nil {
				logWarn("Failed to join channel %s: %v", ch, err)
			} else {
				logInfo("Joined channel: %s", ch)
			}
		}
	}

	logInfo("Completed: %d/%d successful", success, len(phones))
	os.Remove("cache.db")
}
