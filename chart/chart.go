package chart

import (
    "fmt"
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/opts"
    "github.com/zhengow/vngo"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

const BUYICON = "image://data:image/svg+xml;base64,PHN2ZyB0PSIxNjcwMTI1NTQ5ODEwIiBjbGFzcz0iaWNvbiIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHAtaWQ9IjgyOTgiIHdpZHRoPSIyMDAiIGhlaWdodD0iMjAwIj48cGF0aCBkPSJNNTExLjY1MzYxMSA2NC4wNjUxMDhjMCAwLTg5LjQzNDkxNSAzMjYuMjU4NjY5LTMxMC43MzQwOSA1MTcuODc4ODk0IDAgMCAxMDkuMzcwOTY2IDcuMDExNjk0IDI3Ni40NTg0NDktMTA2LjMyMzU2MWwtMC4yNTE3MzMgNDgzLjczMjE5IDY2Ljg0NjQ1NCAwTDU0My45NzI2OTEgNDc1LjkxMjA4NGMwIDAgMTM5LjcwMjc4NiA5Ni42ODUwNCAyNzguNDE4MDgxIDEwNi4wMzE5MThDODIyLjM5MTc5NCA1ODEuOTQ0MDAyIDYxNy4yMDg2NjkgNDA4LjQ5MzYwMSA1MTEuNjUzNjExIDY0LjA2NTEwOHoiIHAtaWQ9IjgyOTkiIGZpbGw9IiNkODFlMDYiPjwvcGF0aD48L3N2Zz4="

const SELLICON = "image://data:image/svg+xml;base64,PHN2ZyB0PSIxNjcwMTI1NTg0NDM3IiBjbGFzcz0iaWNvbiIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHAtaWQ9IjI2MDkiIHdpZHRoPSIyMDAiIGhlaWdodD0iMjAwIj48cGF0aCBkPSJNNTExLjY1NzcwNCA5NTkuMzUxNjA4YzAgMCA4OS40MzQ5MTUtMzI2LjI1ODY2OSAzMTAuNzM0MDktNTE3Ljg3Nzg3MSAwIDAtMTA5LjM3MDk2Ni03LjAxMTY5NC0yNzYuNDU4NDQ5IDEwNi4zMjM1NjFsMC4yNTE3MzMtNDgzLjczMTE2Ni02Ni44NDc0NzcgMCAwIDQ4My40Mzk1MjRjMCAwLTEzOS43MDI3ODYtOTYuNjg1MDQtMjc4LjQxODA4MS0xMDYuMDMxOTE4QzIwMC45MTk1MjEgNDQxLjQ3MzczNyA0MDYuMTAyNjQ2IDYxNC45MjMxMTQgNTExLjY1NzcwNCA5NTkuMzUxNjA4eiIgcC1pZD0iMjYxMCIgZmlsbD0iIzFhZmEyOSI+PC9wYXRoPjwvc3ZnPg=="

func getXData(x []time.Time) []string {
    data := make([]string, len(x))
    for idx, _data := range x {
        data[idx] = _data.Format(vngo.DateFormat)
    }
    return data
}

func getYData(y []float64) []opts.LineData {
    lineData := make([]opts.LineData, len(y))
    for idx, data := range y {
        lineData[idx] = opts.LineData{
            Value:  vngo.RoundTo(data, 2),
            Symbol: "none",
        }
    }
    return lineData
}

func getGlobalOpts() []charts.GlobalOpts {
    titleOpts := charts.WithTitleOpts(opts.Title{
        Title: "PNL",
        Left:  "center",
    })
    initOpts := charts.WithInitializationOpts(opts.Initialization{
        PageTitle: "vngo",
        Width:     "100%",
    })
    toolBoxOpts := charts.WithToolboxOpts(opts.Toolbox{
        Show:  true,
        Right: "5%",
        Feature: &opts.ToolBoxFeature{
            DataZoom: &opts.ToolBoxFeatureDataZoom{
                Show: true,
            },
        },
    })
    toolTipOpts := charts.WithTooltipOpts(opts.Tooltip{
        Show:    true,
        Trigger: "axis",
    })
    dataZoomOpts := charts.WithDataZoomOpts(opts.DataZoom{
        Type: "inside",
    }, opts.DataZoom{
        Type: "slider",
    })
    yAxisOpts := charts.WithYAxisOpts(opts.YAxis{
        Scale: true,
    })
    return []charts.GlobalOpts{titleOpts, initOpts, toolBoxOpts, toolTipOpts, dataZoomOpts, yAxisOpts}
}

func ChartPNL(x []time.Time, y []float64, _filename string) {
    line := charts.NewLine()
    globalOpts := getGlobalOpts()
    line.SetGlobalOptions(globalOpts...)
    lineData := getYData(y)
    line.SetXAxis(getXData(x)).AddSeries("pnl", lineData)
    if _filename != "" {
        f, _ := os.Create(fmt.Sprintf("%s.html", _filename))
        line.Render(f)
    } else {
        f, _ := os.Create("vngo.html")
        line.Render(f)
        http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
            content, _ := ioutil.ReadFile(f.Name())
            w.Write(content)
        })
        fmt.Println("chartPNL: http://localhost:8081")
        http.ListenAndServe(":8081", nil)
    }
}

func getKLineData(bars []vngo.Bar) []opts.KlineData {
    klineData := make([]opts.KlineData, len(bars))
    for idx, bar := range bars {
        klineData[idx] = opts.KlineData{
            Name:  bar.Symbol.Symbol,
            Value: bar.GetKLineData(),
        }
    }
    return klineData
}

func getKLineGlobalOpts() []charts.GlobalOpts {
    titleOpts := charts.WithTitleOpts(opts.Title{
        Title: "KLine",
        Left:  "center",
    })
    initOpts := charts.WithInitializationOpts(opts.Initialization{
        PageTitle: "vngo",
        Width:     "100%",
    })
    toolBoxOpts := charts.WithToolboxOpts(opts.Toolbox{
        Show:  true,
        Right: "5%",
        Feature: &opts.ToolBoxFeature{
            DataZoom: &opts.ToolBoxFeatureDataZoom{
                Show: true,
            },
        },
    })
    toolTipOpts := charts.WithTooltipOpts(opts.Tooltip{
        Show:    true,
        Trigger: "axis",
    })
    dataZoomOpts := charts.WithDataZoomOpts(opts.DataZoom{
        Type: "inside",
    }, opts.DataZoom{
        Type: "slider",
    })
    yAxisOpts := charts.WithYAxisOpts(opts.YAxis{
        Scale: true,
    })
    return []charts.GlobalOpts{titleOpts, initOpts, toolBoxOpts, toolTipOpts, dataZoomOpts, yAxisOpts}
}

func getTradesDataPoints(trades []*vngo.TradeData) []opts.ScatterData {
    scatterData := make([]opts.ScatterData, len(trades))
    for idx, trade := range trades {
        icon := BUYICON
        //rotate := 0
        if trade.IsSell() {
            //rotate = 180
            icon = SELLICON
        }
        scatterData[idx] = opts.ScatterData{
            Name:       trade.Symbol.Symbol,
            Symbol:     icon,
            Value:      [2]interface{}{trade.Datetime.Format(vngo.DateFormat), trade.Price},
            SymbolSize: 50,
            //SymbolRotate: rotate,
        }
    }
    return scatterData
}

func chartTradeScatter(x []time.Time, trades []*vngo.TradeData) *charts.Scatter {
    if len(trades) == 0 {
        return nil
    }
    scatter := charts.NewScatter()
    scatter.SetXAxis(getXData(x)).AddSeries("trade", getTradesDataPoints(trades))
    scatter.SetSeriesOptions(charts.WithLabelOpts(opts.Label{
        Show:      true,
        Formatter: "{@[1]}",
    }))
    return scatter
}

func ChartKLines(x []time.Time, y []vngo.Bar, trades []*vngo.TradeData, _filename string) {
    kline := charts.NewKLine()
    kline.SetGlobalOptions(getKLineGlobalOpts()...)
    kline.SetXAxis(getXData(x)).AddSeries("kline", getKLineData(y))

    scatter := chartTradeScatter(x, trades)
    if scatter != nil {
        kline.Overlap(scatter)
    }
    //kline.SetSeriesOptions(charts.WithMarkPointNameCoordItemOpts(getTradesDataPoints(trades)...), charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
    //    Symbol:     []string{"arrow", "circle"},
    //    SymbolSize: 20,
    //}))
    if _filename != "" {
        f, _ := os.Create(fmt.Sprintf("%s.html", _filename))
        kline.Render(f)
    } else {
        f, _ := os.Create("vngo.html")
        kline.Render(f)
        http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
            content, _ := ioutil.ReadFile(f.Name())
            w.Write(content)
        })
        fmt.Println("chartPNL: http://localhost:8081")
        http.ListenAndServe(":8081", nil)
    }
}
