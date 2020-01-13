package options

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type Configurator interface {
	ConfigPath() string
}

type OptionsCommon struct {
	Config string `short:"c" long:"config" default:"config.ini" no-ini:"true"`

	RestartDelay uint32 `long:"delay"   default:"5"             description:"restart delay seconds on fail"`
}

func (com *OptionsCommon) ConfigPath() string { return com.Config }

type OptionsMQ struct {
	MqHost     string `long:"mqhost"     default:"ip6-localhost" description:"rabbitMQ server ip"`
	MqPort     uint32 `long:"mqport"     default:"5672"          description:"rabbitMQ port"`
	MqUser     string `long:"mquser"     default:"user"          description:"rabbitMQ username"`
	MqPassword string `long:"mqpass"     default:"password"      description:"rabbitMQ user password"`
	//MqQueue    string `long:"mqqueue"    default:"sms"           description:"rabbitMQ queue name"`
	MqTimeout uint32 `long:"mqtimeout"  default:"10"            description:"connection timeout seconds to rabbitMQ"`
}

type OptionsDB struct {
	DbHost     string `long:"dbhost"     default:"ip6-localhost" description:"listening database server ip"`
	DbPort     uint32 `long:"dbport"     default:"5433"          description:"listening database port"`
	DbUser     string `long:"dbuser"     default:"rabbit"        description:"database username"`
	DbPassword string `long:"dbpass"     default:"example"       description:"database user password"`
	DbName     string `long:"dbname"     default:"messages"      description:"database name"`
	DbSsl      string `long:"dbssl"      default:"disable"       description:"database ssl enabled"`
	DbTimeout  uint32 `long:"dbtimeout"  default:"10"            description:"connection timeout seconds to database"`
}

type OptionsWorker struct {
	OptionsCommon
	OptionsMQ
	OptionsDB
}

type OptionsAcceptor struct {
	OptionsCommon
	OptionsMQ

	HttpIP              string `long:"httpIP"            default:"0.0.0.0"  description:"http server ip"`
	HttpPort            uint32 `long:"httpPort"          default:"8080"     description:"http server port"`
	HttpWriteTimeout    uint32 `long:"httpWriteTimeout"  default:"5"        description:"http write timeout"`
	HttpReadTimeout     uint32 `long:"httpReadTimeout"   default:"5"        description:"http read timeout"`
	HttpGracefulTimeout uint32 `long:"httpGrace"         default:"5"        description:"timeout seconds for http server graceful shutdown"`
}

// ReadWorker reads options from file or command line arguments for worker service
func ReadWorker() (OptionsWorker, error) {
	var opts OptionsWorker

	parser, err := parseArgs(&opts)
	if err != nil {
		return opts, fmt.Errorf("parseINI args: %w", err)
	}

	if err := parseINI(&opts, parser); err != nil {
		return opts, fmt.Errorf("parseINI from ini: %w", err)
	}

	return opts, nil
}

// ReadWorker reads options from file or command line arguments for worker service
func ReadAcceptor() (OptionsAcceptor, error) {
	var opts OptionsAcceptor

	parser, err := parseArgs(&opts)
	if err != nil {
		return opts, fmt.Errorf("parseINI args: %w", err)
	}

	if err := parseINI(&opts, parser); err != nil {
		return OptionsAcceptor{}, fmt.Errorf("options.ReadAcceptor: %w", err)
	}

	return opts, nil
}

func parseArgs(opts interface{}) (*flags.Parser, error) {
	parser := flags.NewParser(opts, flags.PassDoubleDash)

	if _, err := flags.ParseArgs(opts, os.Args); err != nil {
		return nil, fmt.Errorf("parse commandline arguments: %w", err)
	}
	return parser, nil
}

func parseINI(opts interface{}, parser *flags.Parser) error {
	cfg, ok := opts.(Configurator)
	if !ok {
		return fmt.Errorf("can't cast opts to Configurator")
	}

	iniPath := cfg.ConfigPath()

	// parse an ini file
	if len(iniPath) > 0 {
		iniParser := flags.NewIniParser(parser)

		if _, err := os.Stat(iniPath); os.IsNotExist(err) {
			// create config file with defaults
			if err := iniParser.WriteFile(iniPath,
				flags.IniIncludeDefaults|flags.IniIncludeComments,
			); err != nil {
				return fmt.Errorf("create default options file: %w", err)
			}
		}

		if err := iniParser.ParseFile(iniPath); err != nil {
			return fmt.Errorf("parseINI options: %w", err)
		}
	}

	return nil
}
