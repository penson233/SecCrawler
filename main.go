package main

import (
	"SecCrawler/bot"
	_ "SecCrawler/bot"
	"SecCrawler/config"
	"SecCrawler/crawler"
	_ "SecCrawler/crawler"
	"SecCrawler/register"
	"SecCrawler/utils"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/robfig/cron"
)

// var cfg = config.Cfg

func init() {
	flag.BoolVar(&config.Test, "test", false, "stop after running once")
	flag.BoolVar(&config.Version, "version", false, "print version info")
	flag.BoolVar(&config.Help, "help", false, "print help info")
	flag.BoolVar(&config.Generate, "init", false, "generate a config file")
	flag.StringVar(&config.ConfigFile, "c", "./config.yml", "the config `file` to be used")
	flag.Usage = usage
}

func usage() {
	fmt.Printf("SecCrawler %s\n\nOptions:\n", config.TAG)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if config.Help {
		flag.Usage()
		return
	}

	if config.Version {
		fmt.Printf("Version: SecCrawler %s\nGithub: %s\nGo Version: %s\nBuild Time: %s\n", config.TAG, config.GITHUB, config.GOVERSION, config.BUILD_TIME)
		return
	}

	config.ConfigInit()
	bot.BotInit()
	crawler.CrawlerInit()

	if config.Test {
		start()
		return
	}

	_cron := cron.New()
	spec := fmt.Sprintf("0 0 %d * * ?", config.Cfg.CronTime)
	err := _cron.AddFunc(spec, start)
	// err := _cron.AddFunc("0 */1 * * * ?", start) //每分钟
	if err != nil {
		log.Fatalf("add cron error: %s\n", err.Error())
	}

	_cron.Start()
	defer _cron.Stop()
	select {}

}

func start() {
	fmt.Printf("%s\n[♥︎] crawler start at %s\n%s\n\n", strings.Repeat("-", 47), utils.CurrentTime(), strings.Repeat("-", 47))

	for crawlerName, crawler := range register.GetCrawlerMap() {
		crawlerResult, err := crawler.Get()
		if err != nil {
			log.Printf("crawl [%s] error: %s\n\n", crawlerName, err.Error())
			continue
		}
		for botName, bot := range register.GetBotMap() {
			err := bot.Send(crawlerResult, crawler.Config().Description)
			if err != nil {
				log.Printf("send [%s] to [%s] error: %s\n", crawlerName, botName, err.Error())
			}
		}
	}

}
