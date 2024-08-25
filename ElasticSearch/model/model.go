package model

type EmployeeDetails struct {
	Id        int32  `json:"id"`
	EmpId     int32  `json:"emp_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
}

type FileDetails struct {
	FileName string `json:"fileName"`
}

type Field struct {
	Key string `json:"key"`
}
