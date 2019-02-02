package handle

import (
    "github.com/ajruckman/lib/err"
    "github.com/kataras/iris"
    "math/rand"
    "strconv"
)

func Testing(ctx iris.Context) {

    var sampleSelectData []string
    for i := 0; i < 50; i++ {
        sampleSelectData = append(sampleSelectData, strconv.Itoa(rand.Intn(20000)))
    }
    ctx.ViewData("sampleSelectData", sampleSelectData)
    ctx.ViewData("sampleTableData", sampleTableData)

    err := ctx.View("/page/testing/testing.gohtml")
    liberr.Err(err)
}
