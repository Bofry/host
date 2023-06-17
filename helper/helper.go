package helper

import "github.com/Bofry/host/internal"

func HostHelper(app *internal.AppModule) *internal.HostHelper {
	return &internal.HostHelper{
		App: app,
	}
}
