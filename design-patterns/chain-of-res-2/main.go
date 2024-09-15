package main

func main() {

	cashier := &Cashier{}

	//Set next for doctor department
	doctor := &Doctor{}
	doctor.setNext(cashier)

	//Set next for reception department
	reception := &Reception{}
	reception.setNext(doctor)

	patient := &Patient{name: "abc"}
	//Patient visiting
	reception.execute(patient)
}
