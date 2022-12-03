package engine

import (
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/opts"
    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/utils"
    "os"
    "time"
)

func getXData(x []time.Time) []string {
    data := make([]string, len(x))
    for idx, _data := range x {
        data[idx] = _data.Format(consts.DateFormat)
    }
    return data
}

func getYData(y []float64) []opts.LineData {
    lineData := make([]opts.LineData, len(y))
    for idx, data := range y {
        lineData[idx] = opts.LineData{
            Value:  utils.RoundTo(data, 2),
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
    return []charts.GlobalOpts{titleOpts, initOpts, toolBoxOpts, toolTipOpts, dataZoomOpts}
}

func chart(x []time.Time, y []float64) {
    line := charts.NewLine()
    globalOpts := getGlobalOpts()
    line.SetGlobalOptions(globalOpts...)
    lineData := getYData(y)
    line.SetXAxis(getXData(x)).AddSeries("pnl", lineData)
    f, _ := os.Create("line.html")
    line.Render(f)
}
