package models

type Summary struct {
	TotalSum   float64 `json:"total_sum"`
	PaidSum    float64 `json:"paid_sum"`
	DueSum     float64 `json:"due_sum"`
	ExpenseSum float64 `json:"expense_sum"`
}

type MonthlyIncome struct {
	Month     int     `json:"month"`
	MonthName string  `json:"month_name"`
	Income    float64 `json:"income"`
}

type CurrentMonthSummary struct {
	TotalSum   float64 `json:"total_sum"`
	PaidSum    float64 `json:"paid_sum"`
	DueSum     float64 `json:"due_sum"`
	ExpenseSum float64 `json:"expense_sum"`
}

func (s *Summary) GetSummary() Summary {
	summary := Summary{}
	db.Raw(`
		SELECT 
			(SELECT SUM(total_amount) FROM invoices) AS total_sum,
			(SELECT SUM(paid_amount) FROM invoices) AS paid_sum,
			(SELECT SUM(due_amount) FROM invoices) AS due_sum,
			(SELECT SUM(amount) FROM expenses) AS expense_sum;
  	`).Scan(&summary)

	return summary
}

func (crs *CurrentMonthSummary) GetCrrMonthSumm() CurrentMonthSummary {
	currSummary := CurrentMonthSummary{}
	db.Raw(`
		SELECT 
			(SELECT IFNULL(SUM(total_amount), 0) FROM invoices WHERE MONTH(issue_date) = MONTH(CURRENT_DATE())) AS total_sum,
			(SELECT IFNULL(SUM(paid_amount), 0) FROM invoices WHERE MONTH(issue_date) = MONTH(CURRENT_DATE())) AS paid_sum,
			(SELECT IFNULL(SUM(due_amount), 0) FROM invoices WHERE MONTH(issue_date) = MONTH(CURRENT_DATE())) AS due_sum,
			(SELECT IFNULL(SUM(amount), 0) FROM expenses WHERE MONTH(date) = MONTH(CURRENT_DATE())) AS expense_sum;
  	`).Scan(&currSummary)

	return currSummary
}

func (mi *MonthlyIncome) GetMonthlyIncomes() []MonthlyIncome {
	monthlyIncomes := []MonthlyIncome{}
	db.Raw(`
		SELECT
			months.month,
			months.month_name,
			COALESCE(SUM(i.total_amount), 0) AS income
		FROM
			(
				SELECT 1 AS month, 'Jan' AS month_name
				UNION SELECT 2, 'Feb'
				UNION SELECT 3, 'Mar'
				UNION SELECT 4, 'Apr'
				UNION SELECT 5, 'May'
				UNION SELECT 6, 'Jun'
				UNION SELECT 7, 'Jul'
				UNION SELECT 8, 'Aug'
				UNION SELECT 9, 'Sep'
				UNION SELECT 10, 'Oct'
				UNION SELECT 11, 'Nov'
				UNION SELECT 12, 'Dec'
			) AS months
		LEFT JOIN
			invoices i ON MONTH(i.issue_date) = months.month
		GROUP BY
			months.month, months.month_name
		ORDER BY
			months.month;
  	`).Scan(&monthlyIncomes)

	return monthlyIncomes
}
