package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBUSER            string
	DBHOST            string
	DBPASS            string
	DBNAME            string
	DBPORT            uint
	CLOUDINARY_CLD    string
	CLOUDINARY_KEY    string
	CLOUDINARY_SECRET string
	CLOUDINARY_FOLDER string
}

func InitConfig() *AppConfig {
	var response = new(AppConfig)
	response = ReadData()
	return response
}

func ReadData() *AppConfig {
	var data = new(AppConfig)

	data = readEnv()

	if data == nil {
		err := godotenv.Load(".env")
		data = readEnv()
		if err != nil || data == nil {
			return nil
		}
	}
	return data
}

func readEnv() *AppConfig {
	var data = new(AppConfig)
	var permit = true

	if val, found := os.LookupEnv("DBUSER"); found {
		data.DBUSER = val
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBPASS"); found {
		data.DBPASS = val
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		data.DBHOST = val
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBPORT"); found {
		cnv, err := strconv.Atoi(val)
		if err != nil {
			permit = false
		}

		data.DBPORT = uint(cnv)
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		data.DBNAME = val
	} else {
		permit = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_CLD"); found {
		data.CLOUDINARY_CLD = val
	} else {
		permit = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_KEY"); found {
		data.CLOUDINARY_KEY = val
	} else {
		permit = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_SECRET"); found {
		data.CLOUDINARY_SECRET = val
	} else {
		permit = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_FOLDER"); found {
		data.CLOUDINARY_FOLDER = val
	} else {
		permit = false
	}

	if !permit {
		return nil
	}

	return data
}
