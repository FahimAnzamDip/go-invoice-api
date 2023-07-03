package handlers

import (
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	summary := (&models.Summary{}).GetSummary()
	currSummary := (&models.CurrentMonthSummary{}).GetCrrMonthSumm()
	monthlyIncomes := (&models.MonthlyIncome{}).GetMonthlyIncomes()

	res := u.Message(true, "")
	res["summary"] = summary
	res["current_month_summary"] = currSummary
	res["monthly_incomes"] = monthlyIncomes

	u.Respond(w, res)
}
