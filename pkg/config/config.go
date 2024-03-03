package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/caarlos0/env/v10"
)

func LoadConfig(schema IConfigSchema) any {
	data := loadConfigFile("mini-oauth2-proxy.json")
	configSchema := parseData(data, schema)
	validateConfig(configSchema)
	return configSchema.CreateConfig()
}

func loadConfigFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return data
}

func parseData(data []byte, schema IConfigSchema) IConfigSchema {
	// IConfigSchemaの実体は構造体のポインタなので、単にschemaを指定してよい
	if err := json.Unmarshal(data, schema); err != nil {
		panic(err)
	}

	// デフォルト値は設定ファイルであり、環境変数は上書きである。
	// なぜなら、環境変数では美しく配列やネストが表現できないため、
	// 環境変数による設定は起動ポートとログレベルといった、
	// アプリケーション自体の動作設定よりも頻繁に変更したくなるであろう設定に限定したいからである。

	// IConfigSchemaの実体は構造体のポインタなので、単にschemaを指定してよい
	if err := env.Parse(schema); err != nil {
		panic(err)
	}

	return schema
}

type IConfigSchema interface {
	Validate() error

	// 任意の設定内容に対応できるようにanyになっている
	// 実際に使うときは適切な型にキャストする必要がある
	CreateConfig() any
}

func validateConfig(config IConfigSchema) {
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}
}
