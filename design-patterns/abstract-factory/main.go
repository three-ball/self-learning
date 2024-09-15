package main

import "fmt"

type Drink interface {
	Drink()
}

type Food interface {
	Eat()
}

type Coffee struct{}

func (c *Coffee) Drink() {
	// drink coffee
	fmt.Println("Drinking Coffee")
}

type Beer struct{}

func (b *Beer) Drink() {
	// drink beer
	fmt.Println("Drinking Beer")
}

type Sandwich struct{}

func (s *Sandwich) Eat() {
	// eat sandwich
	fmt.Println("Eating Sandwich")
}

type Chips struct{}

func (c *Chips) Eat() {
	// eat chips
	fmt.Println("Eating Chips")
}

type VoucherAbstractFactory interface {
	GetDrink() Drink
	GetFood() Food
}

type CoffeeSandwichFactory struct{}

func (c *CoffeeSandwichFactory) GetDrink() Drink {
	return &Coffee{}
}

func (c *CoffeeSandwichFactory) GetFood() Food {
	return &Sandwich{}
}

type BeerChipsFactory struct{}

func (b *BeerChipsFactory) GetDrink() Drink {
	return &Beer{}
}

func (b *BeerChipsFactory) GetFood() Food {
	return &Chips{}
}

func getVoucherFactory(campaignName string) VoucherAbstractFactory {
	if campaignName == "morning" {
		return &CoffeeSandwichFactory{}
	}
	if campaignName == "evening" {
		return &BeerChipsFactory{}
	}
	return nil
}

type Voucher struct {
	drink Drink
	food  Food
}

func GetVoucher(factory VoucherAbstractFactory) Voucher {
	return Voucher{
		drink: factory.GetDrink(),
		food:  factory.GetFood(),
	}
}

func main() {
	morningVoucher := GetVoucher(getVoucherFactory("morning"))
	morningVoucher.drink.Drink()
	morningVoucher.food.Eat()

	eveningVoucher := GetVoucher(getVoucherFactory("evening"))
	eveningVoucher.drink.Drink()
	eveningVoucher.food.Eat()
}
