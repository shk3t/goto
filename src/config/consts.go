package config

import "os"

const MediaPath = "media"
const GotoConfigName = "goto.toml"

var SecretKey = os.Getenv("SECRET_KEY")