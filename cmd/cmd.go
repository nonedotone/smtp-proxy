package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nonedotone/golog"
	"github.com/nonedotone/smtp-proxy/config"
	"github.com/nonedotone/smtp-proxy/mail"
	"github.com/nonedotone/smtp-proxy/mail/gmail"
	"github.com/nonedotone/smtp-proxy/smtp"
	"github.com/nonedotone/smtp-proxy/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:               "smtp-proxy",
		Short:             "Command for mail proxy",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			golog.Log().Level(logLevel)
			if logFile != "" {
				golog.Log().LogFile(logFile).Rolling(logRolling, logInterval)
			}
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version of smtp-proxy",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
		},
	}

	initCmd = &cobra.Command{
		Use:   "init [mail-type(gmail)]",
		Short: "Init mail token config",
		Args:  cobra.MinimumNArgs(1),
		Run:   initRun,
	}
	sendCmd = &cobra.Command{
		Use:   "send [from] [to] [subject] [message]",
		Args:  cobra.ExactArgs(4),
		Short: "Send email",
		Run:   sendRun,
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "serve smtp service",
		Run:   serveRun,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level:debug,info,warn,error")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "log file name (default empty, not write file)")
	rootCmd.PersistentFlags().StringVar(&logRolling, "log-rolling", "", "log rolling mode:time,size(default empty, not rolling)")
	rootCmd.PersistentFlags().Int64Var(&logInterval, "log-interval", 0, "when log-rolling is time,log-interval represent second,\nwhen log-rolling is size,log-interval represent byte")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "./config.json", "path to config")

	rootCmd.AddCommand(versionCmd)

	initCmd.Flags().StringVar(&gmailPermission, "gmail-permission", "", "gmail permission")
	initCmd.Flags().StringVar(&gmailCredentials, "gmail-credentials", "", "path to credentials")
	rootCmd.AddCommand(initCmd)

	sendCmd.Flags().StringVar(&gmailCredentials, "gmail-credentials", "", "path to credentials")
	rootCmd.AddCommand(sendCmd)

	serveCmd.Flags().StringVar(&serveAddress, "address", ":25", "address listen by serve")
	serveCmd.Flags().StringVar(&gmailCredentials, "gmail-credentials", "", "path to credentials")
	rootCmd.AddCommand(serveCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	logLevel    string
	logFile     string
	logRolling  string
	logInterval int64

	configFile       string
	gmailPermission  string
	gmailCredentials string
	serveAddress     string
)

func initRun(_ *cobra.Command, args []string) {
	t := args[0]
	var cfg *config.Config
	var err error
	switch t {
	case config.GmailType:
		credentials, err := gmail.ReadGmailCredentialsOrDefault(gmailCredentials)
		if err != nil {
			golog.Fatalf("read gmail credentials error %v", err)
		}
		cfg, err = gmail.InitConfig(gmail.MailGmailSendScope, credentials)
		if err != nil {
			golog.Fatalf("init gmail token error %v", err)
		}
	default:
		golog.Fatalf("init mail config error %v", err)
	}
	golog.Debugf("init cfg %v", cfg)
	bz, err := json.Marshal(cfg)
	if err != nil {
		golog.Fatalf("marshal config error %v", err)
	}
	err = ioutil.WriteFile(configFile, bz, os.ModePerm)
	if err != nil {
		golog.Fatalf("write file error %v", err)
	}
}

func sendRun(_ *cobra.Command, args []string) {
	from := args[0]
	to := args[1]
	sub := args[2]
	msg := args[3]

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		golog.Errorf("load config %s error %v", configFile, err)
		return
	}
	m, err := mail.NewMail(cfg, gmailCredentials)
	if err != nil {
		golog.Errorf("init mail error %v", err)
		return
	}
	if err := m.Send(from, to, sub, msg); err != nil {
		golog.Errorf("send email error %v", err)
		return
	}
}

func serveRun(_ *cobra.Command, _ []string) {
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		golog.Errorf("load config %s error %v", configFile, err)
		return
	}
	m, err := mail.NewMail(cfg, gmailCredentials)
	if err != nil {
		golog.Errorf("init mail error %v", err)
		return
	}
	h := smtp.NewHandler(serveAddress, m)
	golog.Fatal(h.MailServer())
}
