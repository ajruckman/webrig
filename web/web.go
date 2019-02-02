package web

import (
    "github.com/ajruckman/lib/err"
    "github.com/kataras/iris"

    "github.com/ajruckman/webrig/internal/config"
    "github.com/ajruckman/webrig/web/handle"
)

func Serve() {
    app := iris.Default()

    templates := iris.HTML("./template", ".gohtml")
    templates.Reload(true)
    app.RegisterView(templates)

    app.StaticWeb("/static", "./static")

    app.Get("/", handle.Home)
    app.Get("/testing", handle.Testing)

    err := app.Run(iris.Addr(":" + config.Port))
    if err != iris.ErrServerClosed {
        liberr.Err(err)
    }
}
