# 概要

设计原则：

- 初始化的时候可以panic，记录日志的方法不能触发panic
- 内部错误需要记录日志，无需返回给业务逻辑代码
- 需要限制自身占用的资源

# 使用方法

配置文件

```yaml
logging:
  core: logrus
  level: debug
  formatter: console
  outputs:
    - type: file
      file: /var/log/tsp/ajaxapi.tsp-asset.oss.log
    - type: stdout
    - type: stderr
    - type: rotate_file
      rotate_file:
        file_name: /var/log/tsp/ajaxapi.tsp-asset.oss.log
        max_size: 100         // 100MB
        max_age: 5            // 5days
        max_backups: 3        // 5
    - type: syslog
      syslog:
        address: 127.0.0.1:514,
        facility: local5
        protocol: udp
```

初始化logger

```go
package main

import (
    "gitlab.p1staff.com/tsp/common/log"
    "other_packages..."
)

func main() {
    //...
    config := &log.Config{
        Core:      log.LogrusCore,
        Level:     log.LevelDebug,
        Formatter: log.ConsoleFormater,
        Outputs: []log.Output{
            {
                Type: "stdout",
            },
        },
    }
    logger, _ := log.New(config)

    //...
}
```

记录日志

```go
package any

import (
    "gitlab.p1staff.com/tsp/common/log"
    "other_packages..."
)

func any(ctx context.Context) {
    // ctx := context.WithValue(context.Background(), "trace", map[string]interface{}{
    //    "trace_id": "a_trace_id",
    //    "span_id":  "a_span_id",
    //    "service":  "restapi.sample.tsp",
    //    "hostname": "sample.p1staff.com",
    //    "ip":       "127.0.0.1",
    //})
    logger.WithTraceInCtx(ctx).WithField("k", "v").Errorf("message content, %s", "additional msg")
}

```

# 单元测试

```sh
go test -v -count=1 gitlab.p1staff.com/tsp/common/log/ && \
go test -v -count=1 gitlab.p1staff.com/tsp/common/log/conf && \
go test -v -count=1 gitlab.p1staff.com/tsp/common/log/logger/logrus && \
go test -v -count=1 gitlab.p1staff.com/tsp/common/log/logger/zap
```

注意：如果执行以下命令会报错，原因是并行执行 `log/logger/logrus/logrus_wrapper_test.go` 和 `log/logger/zap/zap_wrapper_test.go` 中的 `TestSuite1` 会产生执行文件冲突

```sh
go test -v -count=1 gitlab.p1staff.com/tsp/common/log/...
```


# 参考

- [日志规范](https://confluence.p1staff.com/pages/viewpage.action?pageId=25072740)
- [tantan-backend-common/log](https://gitlab.p1staff.com/backend/tantan-backend-common/tree/master/log)
