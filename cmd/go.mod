module main

go 1.19

replace (
	api => ../api/ipinfo
	env => ../utils/env
	ip => ../utils/ip
	middleware => ../api/middleware
	password => ../utils/password
	redisDB => ../api/database
	validator => ../utils/validator
)

require (
	api v0.0.0-00010101000000-000000000000
	env v0.0.0-00010101000000-000000000000
)

require (
	github.com/appleboy/gin-jwt/v2 v2.9.0 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20220106215444-fb4bf637b56d // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gin-contrib/cache v1.2.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.8.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/memcachier/mc/v3 v3.0.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/robfig/go-cache v0.0.0-20130306151617-9fc39e0dbf62 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/net v0.3.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	ip v0.0.0-00010101000000-000000000000 // indirect
	middleware v0.0.0-00010101000000-000000000000 // indirect
	password v0.0.0-00010101000000-000000000000 // indirect
	redisDB v0.0.0-00010101000000-000000000000 // indirect
	validator v0.0.0-00010101000000-000000000000 // indirect
)