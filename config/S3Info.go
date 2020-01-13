package config

//import (
//	"gopkg.in/validator.v2"
//)

type S3Info struct {
	Bucket string `validate:"min=3,max=63,regexp=^[a-z0-9.-]*$"`
	Prefix string `validate:"min=1"`
}
