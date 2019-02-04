package handle

import (
    "encoding/json"
    "github.com/ajruckman/lib/err"
    "github.com/kataras/iris"
    "math/rand"
    "strconv"
    "vSummary"
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

func Testing_Data(ctx iris.Context) {
    req, err := vSummary.ParseDatatablesRequest(ctx.FormValues())
    liberr.Err(err)

    var data []sampleTable
    if len(sampleTableData) > req.Start+req.Length {
      data = sampleTableData[req.Start : req.Start+req.Length]
    } else {
      data = sampleTableData[req.Start:]
    }

    var resp = vSummary.DataTablesResponse{
      Draw:            req.Draw,
      RecordsTotal:    len(sampleTableData),
      RecordsFiltered: len(sampleTableData),
      Data:            sampleTableToMap(data),
    }

    _, err = ctx.JSON(resp)
    liberr.Err(err)
}

func sampleTableToMap(table []sampleTable) (res []map[string]interface{}) {
    res = []map[string]interface{}{}

    in, err := json.Marshal(table)
    liberr.Err(err)
    err = json.Unmarshal(in, &res)
    liberr.Err(err)

    return
}
