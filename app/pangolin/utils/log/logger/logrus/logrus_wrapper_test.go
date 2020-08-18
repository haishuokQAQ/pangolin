package logrus

import (
	"context"
	"os"
	"testing"
	"time"

	"gitlab.p1staff.com/tsp/common/log/conf"
	"gitlab.p1staff.com/tsp/common/model"
	"gitlab.p1staff.com/tsp/common/tracing"

	"github.com/stretchr/testify/suite"
	tspctx "gitlab.p1staff.com/tsp/common/context"
)

type TestSuite1 struct {
	suite.Suite

	logFileName    string
	rotateFileName string
}

func (suite *TestSuite1) SetupTest() {
	suite.logFileName = "/tmp/log_file.log"
	suite.rotateFileName = "/tmp/rotate_log_file.log"

	// clear log files
	if _, err := os.Stat(suite.logFileName); err == nil {
		_ = os.Remove(suite.logFileName)
	}

	if _, err := os.Stat(suite.rotateFileName); err == nil {
		_ = os.Remove(suite.rotateFileName)
	}
}

func (suite *TestSuite1) TearDownTest() {
	_ = os.Remove(suite.logFileName)
	_ = os.Remove(suite.rotateFileName)
}

func (suite *TestSuite1) TestNewLogrusWrapper() {
	options := &conf.Config{
		Level:     conf.LevelDebug,
		Formatter: conf.ConsoleFormater,
		Outputs: []conf.Output{
			{
				Type: "stderr",
			},
			{
				Type: "stdout",
			},
			{
				Type: "file",
				File: &suite.logFileName,
			},
			{
				Type: "rotate_file",
				RotateFile: &conf.RotateFile{
					FileName:   suite.rotateFileName,
					MaxSize:    1,
					MaxAge:     3,
					MaxBackups: 10,
					LocalTime:  true,
					Compress:   false,
				},
			},
		},
	}

	wrapper, err := NewLogrusWrapper(options)
	suite.NoErrorf(err, "fail to create wrapper: %v", err)
	defer wrapper.Close()

	wrapper.Flush()
	time.Sleep(1 * time.Second)

	suite.FileExistsf(suite.logFileName, "%v doesn't exist", suite.logFileName)
	suite.FileExistsf(suite.rotateFileName, "%v doesn't exist", suite.rotateFileName)
}

func TestTestSuite1(t *testing.T) {
	suite.Run(t, new(TestSuite1))
}

func TestJSONOutput(t *testing.T) {
	//assert := assert.New(t)

	wrapper := newDebugJSONLoggerWrapper()
	defer wrapper.Close()

	wrapper.Infof("this is a message in json format")
}

func TestSyslogOutput(t *testing.T) {
	// 如何在mac下启动rsyslog参考：
	// https://www.quora.com/How-can-I-enable-syslog-server-on-Mac-OSX-Sierra
	// https://www.rsyslog.com/doc/v8-stable/
	wrapper := newDebugSyslogLoggerWrapper()
	defer wrapper.Close()

	wrapper.Infof("this is a message sent to syslog")

	// /usr/local/var/log/remote.log:
	// 2020-07-24T00:45:09+08:00 192.168.0.106 /Users/{username}/go/src/gitlab.p1staff.com/tsp/common/log/logger/logrus/debug.test[78878]: time="2020-07-24T00:45:09+08:00" level=info msg="this is a message sent to syslog"
}

func TestAllLevelLog(t *testing.T) {
	//assert := assert.New(t)

	wrapper := newDebugTextLoggerWrapper()
	defer wrapper.Close()

	wrapper.Debugf("this is a debug message")
	wrapper.Infof("this is a info message")
	wrapper.Warningf("this is a warning message")
	wrapper.Errorf("this is a error message")
	//wrapper.Fatalf("this is a fatal message %v", time.Now())   // uncomment this line will cause this test case exit with 1

	wrapper.Flush()

	//Outputs:
	//time="2019-07-27T19:11:58+08:00" level=debug msg="this is a debug message"
	//time="2019-07-27T19:11:58+08:00" level=info msg="this is a info message"
	//time="2019-07-27T19:11:58+08:00" level=warning msg="this is a warning message"
	//time="2019-07-27T19:11:58+08:00" level=error msg="this is a error message"
	//time="2019-07-27T19:11:58+08:00" level=info msg="If you see this message and `Flush` is the last logger's method called in your application, it means no log lost."
}

func TestWithField(t *testing.T) {
	//assert := assert.New(t)

	wrapper := newDebugTextLoggerWrapper()
	defer wrapper.Close()

	wrapper.WithField("k1", "v1").
		WithField("k2", "v2").
		Infof("this is a message with two fields k1, k2")

	wrapper.Infof("this message should be without any field")

	newwrapper := wrapper.WithField("k3", "v3").
		WithField("k4", "v4")
	newwrapper.Infof("this message should with two fields k3, k4")
}

func TestWithFields(t *testing.T) {
	//assert := assert.New(t)

	wrapper := newDebugTextLoggerWrapper()
	defer wrapper.Close()

	wrapper.WithFields(map[string]interface{}{"k1": "v1", "k2": "v2"}).
		Infof("this is a message with two fields k1, k2")

	wrapper.Infof("this message should be without any field")

	newwrapper := wrapper.WithFields(map[string]interface{}{"k3": "v3", "k4": "v4"})
	newwrapper.Infof("this message should with two fields k3, k4")
}

func TestWithTraceInCtx(t *testing.T) {
	//assert := assert.New(t)

	wrapper := newDebugTextLoggerWrapper()
	defer wrapper.Close()

	tracing.InitTracer("test")
	_, ctx := tracing.StartSpanFromContext(context.Background(), "op")
	ctx = tspctx.WithService(ctx, "restapi.sample.tsp")
	ctx = tspctx.WithHostname(ctx, "sample.p1staff.com")
	ctx = tspctx.WithIP(ctx, "127.0.0.1")

	ctx = tspctx.WithUser(ctx, model.NewUser("1"))

	wrapper.WithTraceInCtx(ctx).
		Infof("this a message with trace info")

	newwrapper := wrapper.WithTraceInCtx(ctx)
	newwrapper.Infof("this is a another message with trace info")
}

func newDebugTextLoggerWrapper() *LogrusWrapper {
	options := &conf.Config{
		Level:     conf.LevelDebug,
		Formatter: conf.ConsoleFormater,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
		},
	}

	wrapper, _ := NewLogrusWrapper(options)
	return wrapper
}

func newDebugJSONLoggerWrapper() *LogrusWrapper {
	options := &conf.Config{
		Level:     conf.LevelDebug,
		Formatter: conf.JSONFormater,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
		},
	}

	wrapper, _ := NewLogrusWrapper(options)
	return wrapper
}

func newDebugSyslogLoggerWrapper() *LogrusWrapper {
	options := &conf.Config{
		Level:     conf.LevelDebug,
		Formatter: conf.ConsoleFormater,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
			{
				Type: "syslog",
				Syslog: &conf.Syslog{
					Address:  "127.0.0.1:10514",
					Facility: "local5",
					Protocol: "udp",
				},
			},
		},
	}

	wrapper, _ := NewLogrusWrapper(options)
	return wrapper
}
