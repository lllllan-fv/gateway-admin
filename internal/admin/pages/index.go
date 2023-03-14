package pages

import (
	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/smallbox"
)

func GetDashBoard(ctx *context.Context) (types.Panel, error) {

	components := tmpl.Default()
	colComp := components.Col()
	var size = types.SizeMD(3).SM(6).XS(12)

	smallbox1 := smallbox.New().SetColor("blue").SetIcon("ion-ios-gear-outline").SetUrl("/").SetTitle("new users").SetValue("345￥").GetContent()
	smallbox2 := smallbox.New().SetColor("yellow").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("new users").SetValue("80%").GetContent()
	smallbox3 := smallbox.New().SetColor("red").SetIcon("fa-user").SetUrl("/").SetTitle("new users").SetValue("645￥").GetContent()
	smallbox4 := smallbox.New().SetColor("green").SetIcon("ion-ios-cart-outline").SetUrl("/").SetTitle("new users").SetValue("889￥").GetContent()

	col1 := colComp.SetSize(size).SetContent(smallbox1).GetContent()
	col2 := colComp.SetSize(size).SetContent(smallbox2).GetContent()
	col3 := colComp.SetSize(size).SetContent(smallbox3).GetContent()
	col4 := colComp.SetSize(size).SetContent(smallbox4).GetContent()

	row3 := components.Row().SetContent(col1 + col2 + col3 + col4).GetContent()

	return types.Panel{
		Content:     row3,
		Title:       "Dashboard",
		Description: "Dashboard",
	}, nil
}
