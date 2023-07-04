package models

type Summary struct {
	TotalSum     float64 `json:"total_sum"`
	PaidSum      float64 `json:"paid_sum"`
	DueSum       float64 `json:"due_sum"`
	ExpenseSum   float64 `json:"expense_sum"`
	PaidCount    int     `json:"paid_count"`
	UnpaidCount  int     `json:"unpaid_count"`
	PartialCount int     `json:"partial_count"`
	OverdueCount int     `json:"overdue_count"`
}

type MonthlyIncome struct {
	Month     int     `json:"month"`
	MonthName string  `json:"month_name"`
	Income    float64 `json:"income"`
}

type CurrentMonthSummary struct {
	TotalSum     float64 `json:"total_sum"`
	PaidSum      float64 `json:"paid_sum"`
	DueSum       float64 `json:"due_sum"`
	ExpenseSum   float64 `json:"expense_sum"`
	PaidCount    int     `json:"paid_count"`
	UnpaidCount  int     `json:"unpaid_count"`
	PartialCount int     `json:"partial_count"`
	OverdueCount int     `json:"overdue_count"`
}

func (s *Summary) GetSummary() Summary {
	summary := Summary{}
	db.Raw(`
		SELECT
			SUM(total_amount) AS total_sum,
			SUM(paid_amount) AS paid_sum,
			SUM(due_amount) AS due_sum,
			(SELECT SUM(amount) FROM expenses) AS expense_sum,
			SUM(CASE WHEN status = 'Paid' THEN 1 ELSE 0 END) AS paid_count,
			SUM(CASE WHEN status = 'Unpaid' THEN 1 ELSE 0 END) AS unpaid_count,
			SUM(CASE WHEN status = 'Partially Paid' THEN 1 ELSE 0 END) AS partial_count,
			SUM(CASE WHEN status = 'Overdue' THEN 1 ELSE 0 END) AS overdue_count
		FROM
			invoices;
  	`).Scan(&summary)

	return summary
}

func (crs *CurrentMonthSummary) GetCrrMonthSumm() CurrentMonthSummary {
	currSummary := CurrentMonthSummary{}
	db.Raw(`
		SELECT
			COALESCE(SUM(CASE WHEN MONTH(issue_date) = MONTH(CURRENT_DATE()) THEN total_amount END), 0) AS total_sum,
			COALESCE(SUM(CASE WHEN MONTH(issue_date) = MONTH(CURRENT_DATE()) THEN paid_amount END), 0) AS paid_sum,
			COALESCE(SUM(CASE WHEN MONTH(issue_date) = MONTH(CURRENT_DATE()) THEN due_amount END), 0) AS due_sum,
			COALESCE((SELECT SUM(amount) FROM expenses WHERE MONTH(date) = MONTH(CURRENT_DATE())), 0) AS expense_sum,
			COUNT(CASE WHEN MONTH(issue_date) = MONTH(CURRENT_DATE()) AND status = 'Paid' THEN 1 END) AS paid_count,
			COUNT(CASE WHEN MONTH(issue_date) = MONTH(CURRENT_DATE()) AND status = 'Unpaid' THEN 1 END) AS unpaid_count,
			COUNT(CASE WHEN MONTH(issue_date) = MONTH(CURRENT_DATE()) AND status = 'Partially Paid' THEN 1 END) AS partial_count,
			COUNT(CASE WHEN MONTH(issue_date) = MONTH(CURRENT_DATE()) AND status = 'Overdue' THEN 1 END) AS overdue_count
		FROM
			invoices
		WHERE
			MONTH(issue_date) = MONTH(CURRENT_DATE()) OR (SELECT MONTH(date) FROM expenses LIMIT 1) = MONTH(CURRENT_DATE());
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
