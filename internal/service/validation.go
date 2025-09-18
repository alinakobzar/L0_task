package service

import (
	"fmt"

	"L0_task/internal/types"
)

// Функии для валидаци обязательных полей
func ValidateOrder(o types.Order) error {
	if o.OrderUID == "" {
		return fmt.Errorf("order_uid cannot be empty")
	}
	if o.TrackNumber == "" {
		return fmt.Errorf("track_number cannot be empty")
	}
	if o.Entry == "" {
		return fmt.Errorf("entry cannot be empty")
	}
	if o.Locale == "" {
		return fmt.Errorf("locale cannot be empty")
	}
	if o.CustomerID == "" {
		return fmt.Errorf("customer_id cannot be empty")
	}
	if o.DeliveryService == "" {
		return fmt.Errorf("delivery_service cannot be empty")
	}
	if o.ShardKey == "" {
		return fmt.Errorf("shardkey cannot be empty")
	}
	if o.SmID == 0 {
		return fmt.Errorf("sm_id cannot be 0")
	}
	if o.DateCreated.IsZero() {
		return fmt.Errorf("date_created cannot be empty")
	}
	if o.OofShard == "" {
		return fmt.Errorf("oof_shard cannot be empty")
	}
	return nil
}

func ValidateDelivery(d types.Delivery) error {
	if d.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if d.Phone == "" {
		return fmt.Errorf("phone cannot be empty")
	}
	if d.Zip == "" {
		return fmt.Errorf("zip cannot be empty")
	}
	if d.City == "" {
		return fmt.Errorf("city cannot be empty")
	}
	if d.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if d.Region == "" {
		return fmt.Errorf("region cannot be empty")
	}
	if d.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}

func ValidatePayment(p types.Payment) error {
	if p.Transaction == "" {
		return fmt.Errorf("transaction cannot be empty")
	}
	if p.Currency == "" {
		return fmt.Errorf("currency cannot be empty")
	}
	if p.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}
	if p.Amount == 0 {
		return fmt.Errorf("amount cannot be empty")
	}
	if p.PaymentDT == 0 {
		return fmt.Errorf("payment_dt cannot be empty")
	}
	if p.Bank == "" {
		return fmt.Errorf("bank cannot be empty")
	}
	if p.DeliveryCost == 0 {
		return fmt.Errorf("delivery_cost cannot be empty")
	}
	if p.GoodsTotal == 0 {
		return fmt.Errorf("goods_total cannot be empty")
	}
	if p.CustomFee == 0 {
		return fmt.Errorf("custom_fee cannot be empty")
	}
	return nil
}

func ValidateItem(i types.Item) error {
	if i.ChrtID == 0 {
		return fmt.Errorf("chrt_id cannot be empty")
	}
	if i.TrackNumber == "" {
		return fmt.Errorf("track_number cannot be empty")
	}
	if i.Price == 0 {
		return fmt.Errorf("price cannot be empty")
	}
	if i.Rid == "" {
		return fmt.Errorf("rid cannot be empty")
	}
	if i.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if i.Sale == 0 {
		return fmt.Errorf("sale cannot be empty")
	}
	if i.Size == "" {
		return fmt.Errorf("size cannot be empty")
	}
	if i.TotalPrice == 0 {
		return fmt.Errorf("total_price cannot be empty")
	}
	if i.NmID == 0 {
		return fmt.Errorf("nm_id cannot be empty")
	}
	if i.Brand == "" {
		return fmt.Errorf("brand cannot be empty")
	}
	if i.Status == 0 {
		return fmt.Errorf("brand cannot be empty")
	}
	return nil
}
