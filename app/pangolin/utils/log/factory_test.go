package log

import (
	"pangolin/app/pangolin/utils/log/conf"
	"testing"

)

func TestNewZapLogger(t *testing.T) {
	config := &Config{
		Core:      conf.ZapCore,
		Level:     conf.LevelDebug,
		Formatter: conf.ConsoleFormater,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
		},
	}
	logger, _ := New(config)
	logger.WithField("key", "val").Infof("this is test log")
}

func TestNewLogrusLogger(t *testing.T) {
	config := &Config{
		Core:      conf.LogrusCore,
		Level:     conf.LevelDebug,
		Formatter: conf.ConsoleFormater,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
		},
	}
	logger, _ := New(config)
	logger.WithField("key", "val").Infof("this is test log")
}

func TestLogrusAllLevelLog(t *testing.T) {
	//assert := assert.New(t)

	config := &Config{
		Core:      conf.LogrusCore,
		Level:     conf.LevelDebug,
		Formatter: conf.ConsoleFormater,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
		},
	}
	Init(config)
	defer Close()

	Debugf("this is a debug message")
	Infof("this is a info message")
	Warningf("this is a warning message")
	Errorf("this is a errors message")
	//wrapper.Fatalf("this is a fatal message %v", time.Now())   // uncomment this line will cause this test case exit with 1

	Flush()

	//Outputs:
	//time="2019-07-27T19:11:58+08:00" level=debug msg="this is a debug message"
	//time="2019-07-27T19:11:58+08:00" level=info msg="this is a info message"
	//time="2019-07-27T19:11:58+08:00" level=warning msg="this is a warning message"
	//time="2019-07-27T19:11:58+08:00" level=error msg="this is a error message"
	//time="2019-07-27T19:11:58+08:00" level=info msg="If you see this message and `Flush` is the last logger's method called in your application, it means no log lost."
}

func TestZapAllLevelLog(t *testing.T) {
	//assert := assert.New(t)

	config := &conf.Config{
		Core:      conf.ZapCore,
		Level:     conf.LevelDebug,
		Formatter: conf.ConsoleFormater,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
		},
	}
	Init(config)
	defer Close()

	Debugf("this is a debug message")
	Infof("this is a info message")
	Warningf("this is a warning message")
	Errorf("this is a error message")
	//wrapper.Fatalf("this is a fatal message %v", time.Now())   // uncomment this line will cause this test case exit with 1

	Flush()

	//Outputs:
	// 2020-07-23T20:47:04.383+0800    debug   zap/zap_wrapper.go:75   this is a debug message
	// 2020-07-23T20:47:04.383+0800    info    zap/zap_wrapper.go:79   this is a info message
	// 2020-07-23T20:47:04.383+0800    warn    zap/zap_wrapper.go:83   this is a warning message
	// 2020-07-23T20:47:04.383+0800    error   zap/zap_wrapper.go:87   this is a error message
}
