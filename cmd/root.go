package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"tztimer/dbus"
	"tztimer/timer"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	timezone *string
	target   *string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tztimer",
	Short: "A simple timer that sends a notification on D-BUS",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If already detached, just run normally
		if os.Getenv("DAEMONIZED") == "1" {
			return run()
		}

		daemon := exec.Command(os.Args[0], os.Args[1:]...)
		daemon.Env = append(os.Environ(), "DAEMONIZED=1")

		// Detach from terminal
		daemon.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true, // create new session
		}

		// Redirect stdio to avoid terminal attachment
		daemon.Stdin = nil
		daemon.Stdout = nil
		daemon.Stderr = nil

		err := daemon.Start()
		if err != nil {
			return fmt.Errorf("Failed to start detached process: %w", err)
		}

		fmt.Println("Detached with PID:", daemon.Process.Pid)
		return nil
	},
}

func run() error {
	logrus.SetLevel(logrus.InfoLevel)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	bus, err := dbus.NewDbus(ctx)
	if err != nil {
		return err
	}

	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return err
	}

	t, err := timer.New(bus, tz, *target)
	if err != nil {
		return err
	}

	<-t.Start()

	err = bus.Notify(*timezone, fmt.Sprintf("%s", time.Now().In(tz).Format(time.TimeOnly)), 5000)
	if err != nil {
		return err
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	timezone = rootCmd.Flags().StringP("timezone", "z", "Local", "timezone to set alarm based on")
	t := time.Now().Format(time.TimeOnly)
	target = rootCmd.Flags().StringP("time", "t", t, "set alarm for this time hh:mm:ss")
}
