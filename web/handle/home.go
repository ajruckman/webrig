package handle

import (
    "github.com/ajruckman/lib/err"
    "github.com/kataras/iris"
)

func Home(ctx iris.Context) {
    err := ctx.View("/page/home/home.gohtml")
    liberr.Err(err)
}
