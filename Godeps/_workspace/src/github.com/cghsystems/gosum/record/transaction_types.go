package record

const (
	DEBIT  = "DEB"
	CREDIT = "CR"
	FPI    = "FPI"
)

func (record Record) IsFPI() bool {
	return FPI == record.TransactionType
}

func (record Record) IsDebit() bool {
	return DEBIT == record.TransactionType
}

func (record Record) IsCredit() bool {
	return CREDIT == record.TransactionType
}
