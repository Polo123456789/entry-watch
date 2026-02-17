package entry

type AdminUser struct {
	ID            int64
	CondominiumID int64
	FirstName     string
	LastName      string
	Email         string
	Phone         string
	Enabled       bool
	CondoName     string
}
