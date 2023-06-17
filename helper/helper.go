package helper

import (
	. "github.com/Bofry/host/internal"
)

func HostHelper(app *AppModule) *AppHostHelper {
	return &AppHostHelper{
		App: app,
	}
}
