package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thunes/internal/api"
	"github.com/thunes/internal/args"
	"github.com/thunes/internal/config"
	"github.com/thunes/internal/db"
	"github.com/thunes/pkg/errwrap"
)

var (
	_args = args.Args{}
)

func main() {
	fmt.Println("=========================Thunes Start==========================")
	defer fmt.Println("=========================Thunes Exit==========================")

	// panic recover
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("painic: ", p)
		}
	}()

	if err := appInit(); err != nil {
		fmt.Println(err)
		return
	}
	defer appExit()

	// api engine start-up
	if err := appRun(); err != nil {
		fmt.Println(err)
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////
func appInit() (err error) {
	if _args.Parse(); err != nil {
		return errwrap.WithContext(err, "parse args")
	}

	if err = config.Load(_args.CfgPath); err != nil {
		return errwrap.WithContext(err, "load config")
	}

	if err = db.Init(); err != nil {
		return errwrap.WithContext(err, "db init")
	}

	return
}

func appRun() (err error) {
	exitCh := make(chan error)

	go func() {
		exitCh <- api.NewEngine().Run(config.Self.Addr())
	}()

	go func() {
		sig := make(chan os.Signal)

		signal.Notify(sig, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

		s := <-sig

		exitCh <- fmt.Errorf("exit by signal, %s", s)
	}()

	return <-exitCh
}

func appExit() {
	db.Relese()
}
