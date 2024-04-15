package config

import "os"

func InitDirs()  {
    os.Mkdir("media", os.ModePerm)
}