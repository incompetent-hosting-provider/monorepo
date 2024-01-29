package payment

import db_payment "incompetent-hosting-provider/backend/pkg/db/tables"

func GetCurrentCredits(userId string) (int, error) {
	return db_payment.GetUserBalance(userId)
}
