package zap

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

func (suite *TestSuite1) TestWrapperOutPut() {
	options := &conf.Config{
		Level:     conf.LevelDebug,
		Formatter: "console",
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

	wrapper, err := NewZapWrapper(options)
	suite.NoErrorf(err, "fail to create wrapper: %v", err)
	defer wrapper.Close()

	wrapper.Debugf("this is test log")

	wrapper.Flush()
	time.Sleep(1 * time.Second)

	suite.FileExistsf(suite.logFileName, "%v doesn't exist", suite.logFileName)
	suite.FileExistsf(suite.rotateFileName, "%v doesn't exist", suite.rotateFileName)
}

func TestTestSuite1(t *testing.T) {
	suite.Run(t, new(TestSuite1))
}

func newTestZapWrapper(level conf.Level, formatter conf.Formatter) *ZapWrapper {
	options := &conf.Config{
		Level:     level,
		Formatter: formatter,
		Outputs: []conf.Output{
			{
				Type: "stdout",
			},
		},
	}

	wrapper, _ := NewZapWrapper(options)
	return wrapper
}

func TestWithJson(t *testing.T) {
	wrapper := newTestZapWrapper(conf.LevelDebug, conf.ConsoleFormater)
	wrapper.Debugf("this is test json log")
	wrapper.Flush()
}

func TestWithConcole(t *testing.T) {
	wrapper := newTestZapWrapper(conf.LevelDebug, conf.ConsoleFormater)
	wrapper.Debugf("this is test console log")
	wrapper.Flush()
}

func TestWithLevel(t *testing.T) {
	wrapper := newTestZapWrapper(conf.LevelInfo, conf.ConsoleFormater)
	wrapper.Debugf("this is test debug log")
	wrapper.Infof("this is test info log")
	wrapper.Warningf("this is test warn log")
	wrapper.Errorf("this is test error log")
	wrapper.Flush()

	// Output:
	// 2019-07-29T15:26:42.752+0800    info    zap/zap_wrapper.go:131  this is test info log
	// 2019-07-29T15:26:42.752+0800    warn    zap/zap_wrapper.go:135  this is test warn log
	// 2019-07-29T15:26:42.752+0800    error   zap/zap_wrapper.go:139  this is test error log

	wrapper = newTestZapWrapper(conf.LevelInfo, conf.JSONFormater)
	wrapper.Debugf("this is test debug log")
	wrapper.Infof("this is test info log")
	wrapper.Warningf("this is test warn log")
	wrapper.Errorf("this is test error log")
	wrapper.Flush()

	// Output:
	// {"level":"info","time":"2019-07-29T15:28:04.821+0800","caller":"zap/zap_wrapper.go:131","msg":"this is test info log"}
	// {"level":"warn","time":"2019-07-29T15:28:04.821+0800","caller":"zap/zap_wrapper.go:135","msg":"this is test warn log"}
	// {"level":"error","time":"2019-07-29T15:28:04.821+0800","caller":"zap/zap_wrapper.go:139","msg":"this is test error log","stacktrace":"gitlab.p1staff.com/tsp/time/log/logger/zap.(*ZapWrapper).Errorf\n\t/Users/huangkai/go/src/gitlab.p1staff.com/tsp/time/log/logger/zap/zap_wrapper.go:139\ngitlab.p1staff.com/tsp/time/log/logger/zap.TestWithLevel\n\t/Users/huangkai/go/src/gitlab.p1staff.com/tsp/time/log/logger/zap/zap_wrapper_test.go:128\ntesting.tRunner\n\t/usr/local/go/src/testing/testing.go:865"}
}

func TestWithField(t *testing.T) {
	wrapper := newTestZapWrapper(conf.LevelDebug, conf.JSONFormater)
	wrapper.WithField("testField", "testValue").WithField("testField1", "testValue1").Debugf("this is test debug log")
	wrapper.Flush()

	wrapper = newTestZapWrapper(conf.LevelDebug, conf.ConsoleFormater)
	wrapper.WithField("testField", "testValue").WithField("testField1", "testValue1").Debugf("this is test debug log")
	wrapper.Flush()
}

func TestWithFields(t *testing.T) {
	wrapper := newTestZapWrapper(conf.LevelDebug, conf.JSONFormater)
	wrapper.WithField("testField3", "testValue3").WithFields(map[string]interface{}{"testField": "testValue", "testField1": "testValue1"}).Debugf("this is test debug log")
	wrapper.Flush()

	wrapper = newTestZapWrapper(conf.LevelDebug, conf.ConsoleFormater)
	wrapper.WithField("testField3", "testValue3").WithFields(map[string]interface{}{"testField": "testValue", "testField1": "testValue1"}).Debugf("this is test debug log")
	wrapper.Flush()
}

func TestWithTraceInCtx(t *testing.T) {
	tracing.InitTracer("test")
	_, ctx := tracing.StartSpanFromContext(context.Background(), "op")
	ctx = tspctx.WithService(ctx, "restapi.sample.tsp")
	ctx = tspctx.WithHostname(ctx, "sample.p1staff.com")
	ctx = tspctx.WithIP(ctx, "127.0.0.1")
	ctx = tspctx.WithUser(ctx, model.NewUser("1"))

	wrapper := newTestZapWrapper(conf.LevelDebug, conf.JSONFormater)
	wrapper.WithTraceInCtx(ctx).Debugf("this is test trace log")
	wrapper.Flush()
}

func TestSyslogOutput(t *testing.T) {
	// 如何在mac下启动rsyslog参考：
	// https://www.quora.com/How-can-I-enable-syslog-server-on-Mac-OSX-Sierra
	// https://www.rsyslog.com/doc/v8-stable/
	wrapper := newDebugSyslogLoggerWrapper()
	defer wrapper.Close()

	wrapper.Infof("this is a message sent to syslog using ZapWrapper")

	// /usr/local/var/log/remote.log:
	// 2020-07-24T01:13:48+08:00 192.168.0.106 /Users/{username}/go/src/gitlab.p1staff.com/tsp/common/log/logger/zap/debug.test[83288]: 2020-07-24T01:13:48.142+0800#011info#011zap/zap_wrapper.go:90#011this is a message sent to syslog using ZapWrapper
}

func newDebugSyslogLoggerWrapper() *ZapWrapper {
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

	wrapper, _ := NewZapWrapper(options)
	return wrapper
}
