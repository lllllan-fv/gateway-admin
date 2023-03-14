package pages

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/smallbox"
	"github.com/GoAdminGroup/themes/sword/components/description"
	"github.com/GoAdminGroup/themes/sword/components/progress_group"
)

func GetDashBoard(ctx *context.Context) (types.Panel, error) {

	components := tmpl.Default()
	colComp := components.Col()

	// ===========================================================================================================================================

	smallbox1 := smallbox.New().SetColor("green").SetIcon("fa-user").SetUrl("/").SetTitle("租户数").SetValue("5").GetContent()
	smallbox2 := smallbox.New().SetColor("blue").SetIcon("ion-ios-gear-outline").SetUrl("/").SetTitle("服务数").SetValue("10").GetContent()
	smallbox3 := smallbox.New().SetColor("yellow").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("当日请求量").SetValue("1120").GetContent()
	smallbox4 := smallbox.New().SetColor("red").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("当前 QPS").SetValue("5").GetContent()

	var size = types.SizeMD(3).SM(6).XS(12)
	col1 := colComp.SetSize(size).SetContent(smallbox1).GetContent()
	col2 := colComp.SetSize(size).SetContent(smallbox2).GetContent()
	col3 := colComp.SetSize(size).SetContent(smallbox3).GetContent()
	col4 := colComp.SetSize(size).SetContent(smallbox4).GetContent()

	row1 := components.Row().SetContent(col1 + col2 + col3 + col4).GetContent()

	// ===========================================================================================================================================

	line := chartjs.Line()

	lineChart := line.
		SetID("salechart").
		SetHeight(180).
		SetTitle("Sales: 1 Jan, 2019 - 30 Jul, 2019").
		SetLabels([]string{"January", "February", "March", "April", "May", "June", "July"}).
		AddDataSet("Electronics").
		DSData([]float64{65, 59, 80, 81, 56, 55, 40}).
		DSFill(false).
		DSBorderColor("rgb(210, 214, 222)").
		DSLineTension(0.1).
		AddDataSet("Digital Goods").
		DSData([]float64{28, 48, 40, 19, 86, 27, 90}).
		DSFill(false).
		DSBorderColor("rgba(60,141,188,1)").
		DSLineTension(0.1).
		GetContent()

	title := `<p class="text-center"><strong>Goal Completion</strong></p>`
	progressGroup := progress_group.New().
		SetTitle("Add Products to Cart").
		SetColor("#76b2d4").
		SetDenominator(200).
		SetMolecular(160).
		SetPercent(80).
		GetContent()

	progressGroup1 := progress_group.New().
		SetTitle("Complete Purchase").
		SetColor("#f17c6e").
		SetDenominator(400).
		SetMolecular(310).
		SetPercent(80).
		GetContent()

	progressGroup2 := progress_group.New().
		SetTitle("Visit Premium Page").
		SetColor("#ace0ae").
		SetDenominator(800).
		SetMolecular(490).
		SetPercent(80).
		GetContent()

	progressGroup3 := progress_group.New().
		SetTitle("Send Inquiries").
		SetColor("#fdd698").
		SetDenominator(500).
		SetMolecular(250).
		SetPercent(50).
		GetContent()

	boxInternalCol1 := colComp.SetContent(lineChart).SetSize(types.SizeMD(8)).GetContent()
	boxInternalCol2 := colComp.
		SetContent(template.HTML(title) + progressGroup + progressGroup1 + progressGroup2 + progressGroup3).
		SetSize(types.SizeMD(4)).
		GetContent()

	boxInternalRow := components.Row().SetContent(boxInternalCol1 + boxInternalCol2).GetContent()

	description1 := description.New().
		SetPercent("17").
		SetNumber("¥140,100").
		SetTitle("TOTAL REVENUE").
		SetArrow("up").
		SetColor("green").
		SetBorder("right").
		GetContent()

	description2 := description.New().
		SetPercent("2").
		SetNumber("440,560").
		SetTitle("TOTAL REVENUE").
		SetArrow("down").
		SetColor("red").
		SetBorder("right").
		GetContent()

	description3 := description.New().
		SetPercent("12").
		SetNumber("¥140,050").
		SetTitle("TOTAL REVENUE").
		SetArrow("up").
		SetColor("green").
		SetBorder("right").
		GetContent()

	description4 := description.New().
		SetPercent("1").
		SetNumber("30943").
		SetTitle("TOTAL REVENUE").
		SetArrow("up").
		SetColor("green").
		GetContent()

	size2 := types.SizeSM(3).XS(6)
	boxInternalCol3 := colComp.SetContent(description1).SetSize(size2).GetContent()
	boxInternalCol4 := colComp.SetContent(description2).SetSize(size2).GetContent()
	boxInternalCol5 := colComp.SetContent(description3).SetSize(size2).GetContent()
	boxInternalCol6 := colComp.SetContent(description4).SetSize(size2).GetContent()

	boxInternalRow2 := components.Row().SetContent(boxInternalCol3 + boxInternalCol4 + boxInternalCol5 + boxInternalCol6).GetContent()

	box := components.Box().WithHeadBorder().SetHeader("Monthly Recap Report").
		SetBody(boxInternalRow).
		SetFooter(boxInternalRow2).
		GetContent()

	boxcol := colComp.SetContent(box).SetSize(types.SizeMD(12)).GetContent()
	row2 := components.Row().SetContent(boxcol).GetContent()

	return types.Panel{
		Content:     row1 + row2,
		Title:       "Dashboard",
		Description: "Dashboard",
	}, nil
}
