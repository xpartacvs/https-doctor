package config

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type config struct {
	hosts          []string
	logLevel       zerolog.Level
	schedule       string
	dishookBotMsg  string
	dishookBotName string
}

type Config interface {
	Hosts() []string
	ZerologLevel() zerolog.Level
	Schedule() string
	DishookBotMessage() string
	DishookBotName() string
}

var (
	cfg  *config
	once sync.Once
)

func (c *config) Hosts() []string {
	return c.hosts
}

func (c *config) ZerologLevel() zerolog.Level {
	return c.logLevel
}

func (c *config) Schedule() string {
	return c.schedule
}

func (c *config) DishookBotMessage() string {
	return c.dishookBotMsg
}

func (c *config) DishookBotName() string {
	return c.DishookBotName()
}

func load() *config {
	fang := viper.New()

	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()

	fang.SetConfigName("https-doctor")
	fang.SetConfigType("yml")
	fang.AddConfigPath(".")

	value, available := os.LookupEnv("CONFIG_LOCATION")
	if available {
		fang.AddConfigPath(value)
	}

	_ = fang.ReadInConfig()

	return &config{
		hosts:          splitCSV(fang.GetString("hosts")),
		logLevel:       setLogLevel(fang.GetString("loglevel")),
		schedule:       setDefaultString(fang.GetString("schedule"), "0 0 * * *", true),
		dishookBotMsg:  setDefaultString(fang.GetString("dishook.bot.message"), "Your HTTPS health monitoring result", true),
		dishookBotName: setDefaultString(fang.GetString("dishook.bot.name"), "HTTPS Doctor", true),
	}
}

func splitCSV(s string) []string {
	strTrimmed := strings.TrimFunc(
		s,
		func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		},
	)
	comaToSpace := strings.NewReplacer(",", " ")
	strReplacedComa := comaToSpace.Replace(strTrimmed)
	rgxSpaceSplit := regexp.MustCompile(`\s+`)
	return rgxSpaceSplit.Split(strReplacedComa, -1)
}

func setLogLevel(l string) zerolog.Level {
	switch l {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.Disabled
	}
}

func setDefaultString(value, fallback string, trimSpace bool) string {
	if trimSpace {
		value = strings.TrimSpace(value)
	}
	if len(value) <= 0 {
		return fallback
	}
	return value
}

func Get() Config {
	once.Do(func() {
		cfg = load()
	})
	return cfg
}
