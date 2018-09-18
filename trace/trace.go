package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"encoding/json"
)

func saveAccount(stub shim.ChaincodeStubInterface, account Account) bool {
	acc, err := json.Marshal(account)
	if err != nil {
		return false
	}

	err = stub.PutState(account.CardNo, acc)
	if err != nil {
		return false
	}
	return true
}

func GetAccountByNo(stub shim.ChaincodeStubInterface, cardNo string) (Account, bool) {
	var account Account
	result, err := stub.GetState(cardNo)
	if err != nil {
		return account, false
	}

	err = json.Unmarshal(result, &account)
	if err != nil {
		return account, false
	}

	return account, true
}

// 实现贷款功能
// -c '{"Args":["loan", "身份证号码", "银行名称", "贷款金额"]}'
func loan(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	am, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("给定的贷款金额错误")
	}

	bank := Bank{
		BankName: args[1],
		Flag:Bank_Flag_Loan,
		Amount:am,
		StartDate:"20100901",
		EndDate:"20101201",
	}

	account := Account{
		CardNo:args[0],
		Aname:"jack",
		Age:29,
		Gender:"男",
		Mobil:"6234567",
		Bank:bank,
	}

	bl := saveAccount(stub, account)
	if !bl {
		return shim.Error("保存贷款记录失败")
	}

	return shim.Success([]byte("贷款成功"))
}

// 实现还款功能
// -c '{"Args":["loan", "身份证号码", "银行名称", "还款金额"]}'
func repayment(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	am, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("给定的贷款金额错误")
	}

	bank := Bank{
		BankName: args[1],
		Flag:Bank_Flag_Repayment,
		Amount:am,
		StartDate:"20101001",
		EndDate:"20101201",
	}

	account := Account{
		CardNo:args[0],
		Aname:"jack",
		Age:29,
		Gender:"男",
		Mobil:"6234567",
		Bank:bank,
	}

	bl := saveAccount(stub, account)
	if !bl {
		return shim.Error("保存还款记录失败")
	}

	return shim.Success([]byte("此次还款成功"))

}

// 根据账户身份证号码查询相应的信息(包含该账户所有的历史记录信息)
// -c '{"Args":["queryAccountByCardNo", "身份证号码]}'
func queryAccountByCardNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("必须且只能指定要查询的账户信息的身份证号码")
	}

	account, bl := GetAccountByNo(stub, args[0])
	if !bl{
		return shim.Error("根据指定的身份证号码查询信息时发生错误")
	}

	// 查询历史记录信息
	accIterator, err := stub.GetHistoryForKey(account.CardNo)
	if err != nil {
		return shim.Error("查询历史记录信息时发生错误")
	}
	defer accIterator.Close()

	// 处理查询到的历史记录信息迭代器对象
	var historys []HistoryItem
	var acc Account
	for accIterator.HasNext() {
		// 依次获取迭代器中的元素
		hisData, err := accIterator.Next()
		if err != nil {
			return shim.Error("处理迭代器对象时发生错误")
		}

		var hisItem HistoryItem
		hisItem.TxId = hisData.TxId		// 获取此条交易的交易编号
		err = json.Unmarshal(hisData.Value, &acc)	// 获取此条交易的状态信息
		if err != nil {
			return shim.Error("反序列化历史状态时发生错误")
		}
		// 处理当前记录状态为nil的情况
		if hisData.Value == nil {
			var empty Account
			hisItem.Account = empty
		}else{
			hisItem.Account = acc
		}

		// 将当前处理完毕的历史状态保存至数组中
		historys = append(historys, hisItem)
	}

	account.Historys = historys

	accByte, err := json.Marshal(account)
	if err != nil {
		return shim.Error("将账户信息序列化时发生错误")
	}

	return shim.Success(accByte)
}
