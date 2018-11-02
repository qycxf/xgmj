//
// Author: leafsoar
// Date: 2016-05-05 18:35:18
//

// 游戏大厅服务器 (http)

package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "expvar"
	_ "net/http/pprof"

	_ "github.com/lib/pq"
	"qianuuu.com/ahmj/internal/config"
	"qianuuu.com/ahmj/internal/game"
	"qianuuu.com/lib/logs"
)

func handleSignals() {
	// 添加性能分析工具
	pfport := os.Getenv("GO_PPROF")
	if pfport != "" {
		// import _ "expvar"
		// import _ "net/http/pprof"
		// http://localhost:6060/debug/vars
		// http://localhost:6060/debug/pprof
		// http://localhost:6060/debug/pprof/profile
		go func() {
			logs.Info("[main] start pprof port %s ...", pfport)
			err := http.ListenAndServe(":"+pfport, nil)
			if err != nil {
				logs.Error("[main] pprof " + err.Error())
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func main() {
	// runtime.GOMAXPROCS(1)

	// // 添加性能分析工具
	// go func() {
	// 	paymentProfile := pprof.NewProfile("profile")
	// 	_ = paymentProfile
	// 	http.ListenAndServe("localhost:"+config.Opts.ProfilePort, nil)
	// 	http.DefaultServeMux.Handle("/debug/pprof/profile", pprofHTTP.Handler("profile"))
	// }()

	// 设置时区
	time.Local = time.FixedZone("CST", 8*3600)

	// 读取配置文件
	if err := config.ParseToml("config.toml"); err != nil {
		logs.Info("配置文件读取失败: %s", err.Error())
		return
	}
	config.ReadJson()
	// 日志设置
	if len(config.Opts().LogPath) > 0 {

		logStr := "debug"
		// if config.Opts.CloseTest {
		// 	logStr = "release"
		// }
		log, err := logs.New(logStr, config.Opts().LogPath)
		if err == nil {
			logs.Export(log)
			defer func() {
				log.Close()
			}()
		} else {
			logs.Error(err.Error())
		}
	}
	defer func() {
		logs.Close()
	}()

	// login.Main()
	g := game.NewGame()
	go g.Serve()

	handleSignals()
}
