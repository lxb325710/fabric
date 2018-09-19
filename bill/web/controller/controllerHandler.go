package controller

import (
	"net/http"
	"github.com/kongyixueyuan.com/bill/service"
	"encoding/json"
	"fmt"
)

type Application struct {
	Fabric *service.FabricSetupService
}

var cuser User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request)  {
	response(w, r, "login.html", nil)
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request)  {
	userName := r.FormValue("userName")
	password := r.FormValue("password")

	var flag = false
	for _, user := range Users {
		if userName == user.UserName && password == user.Password {
			cuser = user
			flag = true
			break
		}
	}

	if flag{
		app.FindBills(w, r)
	}else {
		data := &struct {
			Name string
			Flag bool
		}{
			Name:userName,
			Flag:true,
		}
		response(w, r, "login.html", data)
	}

}

// 显示发布票据页面
func (app *Application) IssueView(w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		Flag bool
		Msg string
		Cuser User
	}{
		Flag:false,
		Msg:"",
		Cuser:cuser,
	}
	response(w, r, "issue.html", data)
}

// 发布新票据
func (app *Application) SaveBill(w http.ResponseWriter, r *http.Request)  {
	bill := service.Bill{
		BillInfoID: r.FormValue("BillInfoID"),
		BillInfoType:r.FormValue("BillInfoType"),
		BillInfoAmt:r.FormValue("BillInfoAmt"),

		DrwrAcct: r.FormValue("DrwrAcct"),
		DrwrCmID: r.FormValue("DrwrCmID"),

		AccptrAcct: r.FormValue("AccptrAcct"),
		AccptrCmID: r.FormValue("AccptrCmID"),

		PyeeAcct: r.FormValue("PyeeAcct"),
		PyeeCmID: r.FormValue("PyeeCmID"),

		HoldrAcct: r.FormValue("HoldrAcct"),
		HoldrCmID: r.FormValue("HoldrCmID"),
	}

	transactionId, err := app.Fabric.SaveBill(bill)
	var msg string
	if err != nil {
		msg = "票据发布失败: " + err.Error()
	}else{
		msg = "票据发布成功: " + transactionId
	}

	data := &struct {
		Msg string
		Flag bool
		Cuser User
	}{
		Msg: msg,
		Flag:true,
		Cuser:cuser,
	}
	response(w, r, "issue.html", data)

}

// 查询我的票据列表
func (app *Application) FindBills(w http.ResponseWriter, r *http.Request)  {
	result, err := app.Fabric.QueryBills(cuser.CmId)
	if err != nil{
		fmt.Println("查询票据列表错误: %v", err)
	}
	var bills []service.Bill
	json.Unmarshal(result, &bills)

	data := &struct {
		Bills []service.Bill
		Cuser User
	}{
		Bills:bills,
		Cuser:cuser,
	}

	response(w, r, "bills.html", data)
}



