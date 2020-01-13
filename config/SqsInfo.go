package config

type SqsInfo struct {
	SqsUrl string `validate:"min=10"`
}
