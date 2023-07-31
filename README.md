# Boot Package README

在monorepo中，當package互相引用越來越複雜後，config管理是一個挑戰：注入所有參數可能導致資訊洩漏，缺少參數則可能讓程式不正確運作。

這個package可以幫助您管理程式的生命週期和所需參數，藉由在init期間向此package註冊參數，我們可以在程式啟動時準確獲知有哪些參數是必備的，並嘗試使用參數先行初始化。

這個package也引用了`github.com/urfave/cli/v2`作為Args解析的工具

## Installation

```
go get github.com/sky-mirror/boot
```

## Usage

### For Singleton package

這個類型的package特徵是藉由sync.Once或init來初始化物件，並提供一個`Default()`或`Get()`函數來回傳同一個物件。

您需要定義一個`config struct`並進行註冊

```go
type config struct {
    Path string
}

var defaultConfig config

func init() {
    boot.Register(&defaultConfig)
}

func (cfg *config) CliFlags() []cli.Flag {
    var flags []cli.Flag
    flags = append(flags, &cli.StringFlag{
      Name: "path",
      Destination: &cfg.Path,
    })
    return flags
}
```

`examples/monorepo/logger/cli-flag.go` 可作為範例參考

### For API wrapper package

這類的package通常對API進行了一個簡單的封裝，並且需要填入伺服器位址，在使用上有可能需要多個實例，僅有config不同。

由於command line options的名稱必須是獨一無二的，我們建議將註冊工作移動到caller side

在callee side

```
func NewConfig(prefix string) *Config {
    cfg := &Config{...}
    boot.Register(cfg)
    return cfg
}

func NewClient(*Config) *Client {
}
```

在caller side

```
var twConfig = api.NewConfig("tw")
var usConfig = api.NewConfig("us")

func Start(ctx context.Context) {
    twClient := api.NewClient(twConfig)
    usClient := api.NewClient(usConfig)
}
```

`examples/monorepo/slack/cli-flag.go` 可作為範例參考

### Before / After hook

如果您註冊的物件實作了`Beforer() / Afterer()`介面，那麼它也會被自動呼叫，呼叫的順序是根據init執行的順序(After則相反)

當它回傳error時程式則會終止執行，所以您可以用來檢測參數正確性並進行初始化工作

### Main function

在`boot.App`中會進行上述的參數處理和生命週期管理，同時負責監聽是否有SIGINT、SIGTERM訊號來發起程式終止流程

您只需要實作Main()內部流程並觀測ctx是否處於結束中的狀態。

```
func Main(ctx context.Context, c *cli.Context) {
    <-ctx.Done()
}

func main() {
    app := boot.App{
        Main: Main,
    }

    app.Run()
}
```



### Results

執行 `go run examples/monorepo/cmd/api/main.go` 可以觀察結果


```
NAME:
   main - A new cli application

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-enable-file                (default: false) [$LOG_ENABLE_FILE]
   --log-file-name value            filename prefix of log file (default: "main") [$LOG_FILE_NAME]
   --log-file-dir value             path of log file (default: "/tmp") [$LOG_FILE_DIR]
   --alert-slack-webhook-url value
   --alert-slack-channel value
   --info-slack-webhook-url value
   --info-slack-channel value
   --help, -h                       show help
2023/07/31 12:54:09 running *logger.config
2023/07/31 12:54:09 Required flags "alert-slack-webhook-url, alert-slack-channel, info-slack-webhook-url, info-slack-channel" not set
```
