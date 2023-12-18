package midtrans

import (
	"BE-hi-SPEC/config"
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

func MidtransCreateToken(orderID int, TotalPrice int) *snap.Response {
	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(config.InitConfig().MIDTRANS_KEY, midtrans.Sandbox)
	id := strconv.Itoa(orderID)
	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "YOUR-ORDER-ID-" + id,
			GrossAmt: int64(TotalPrice),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, _ := s.CreateTransaction(req)
	fmt.Println("Response :", snapResp)
	return snapResp
}

func MidtransStatus(orderID int) (Status string) {
	var c = coreapi.Client{}
	c.New(config.InitConfig().MIDTRANS_KEY, midtrans.Sandbox)
	id := strconv.Itoa(orderID)
	orderId := "YOUR-ORDER-ID-" + id

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := c.CheckTransaction(orderId)
	if e != nil {
		status := "Pending"
		return status
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					status := "Challange"
					return status
				} else if transactionStatusResp.FraudStatus == "accept" {
					status := "Accept"
					return status
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				status := "Success"
				return status
			} else if transactionStatusResp.TransactionStatus == "deny" {
				status := "Deny"
				return status
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				status := "Canceled"
				return status
			} else if transactionStatusResp.TransactionStatus == "pending" {
				status := "Pending"
				return status
			}
		}
	}

	status := "Pending"
	return status
}
