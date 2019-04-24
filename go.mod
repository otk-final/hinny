module github.com/otk-final/hinny

go 1.12

replace google.golang.org/appengine => github.com/golang/appengine v1.5.0

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20190422233926-fe54fb35175b

replace golang.org/x/net => github.com/golang/net v0.0.0-20190420063019-afa5a82059c6

replace golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190402181905-9f3314589c9a

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190422165155-953cdadca894

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

replace golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190422183909-d864b10871cd

replace golang.org/x/vgo => github.com/golang/vgo v0.0.0-20180912184537-9d567625acf4

replace golang.org/x/text => github.com/golang/text v0.3.1-0.20190410012825-f4905fbd45b6

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/core v0.6.2
	github.com/go-xorm/xorm v0.7.1
	github.com/gorilla/mux v1.7.1
	github.com/otk-final/lets-go v0.0.0-20190424032848-f5f17058f52e
	github.com/robertkrimen/otto v0.0.0-20180617131154-15f95af6e78d
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/spf13/viper v1.3.2
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)
