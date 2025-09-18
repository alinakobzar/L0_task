package db

import (
	"L0_task/internal/types"
	"database/sql"
)

// Функция SaveOrder сохраняет заказ в бд
func SaveOrder(db *sql.DB, o types.Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// orders
	_, err = tx.Exec(`
        INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
    `, o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard)
	if err != nil {
		return err
	}

	// delivery
	_, err = tx.Exec(`
        INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8)
    `, o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		return err
	}

	// payment
	_, err = tx.Exec(`
        INSERT INTO payment(order_uid, transaction, request_id, currency, provider, amount,
            payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
    `, o.OrderUID, o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency,
		o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDT, o.Payment.Bank,
		o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee)
	if err != nil {
		return err
	}

	// items
	for _, item := range o.Items {
		_, err = tx.Exec(`
            INSERT INTO items(order_uid, chrt_id, track_number, price, rid, name, sale, size,
                total_price, nm_id, brand, status)
            VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
        `, o.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

//Функия LoadSingleOrder загружает один заказ из бд
func LoadSingleOrder(db *sql.DB, orderUID string) (types.Order, error) {
	var o types.Order

	// orders
	err := db.QueryRow(`SELECT order_uid, track_number, entry, locale, internal_signature,
        customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders WHERE order_uid=$1`, orderUID).
		Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
			&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard)
	if err != nil {
		return o, err
	}

	// delivery
	var d types.Delivery
	err = db.QueryRow(`SELECT name, phone, zip, city, address, region, email
        FROM delivery WHERE order_uid=$1`, orderUID).
		Scan(&d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email)
	if err != nil {
		return o, err
	}
	o.Delivery = d

	// payment
	var p types.Payment
	err = db.QueryRow(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank,
        delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid=$1`, orderUID).
		Scan(&p.Transaction, &p.RequestID, &p.Currency, &p.Provider, &p.Amount, &p.PaymentDT,
			&p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee)
	if err != nil {
		return o, err
	}
	o.Payment = p

	// items
	rows, err := db.Query(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price,
        nm_id, brand, status FROM items WHERE order_uid=$1`, orderUID)
	if err != nil {
		return o, err
	}
	defer rows.Close()

	var items []types.Item
	for rows.Next() {
		var it types.Item
		if err := rows.Scan(&it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid, &it.Name,
			&it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status); err != nil {
			return o, err
		}
		items = append(items, it)
	}
	o.Items = items

	return o, nil
}

//Функция LoadAllOrders загружает все заказы из бд
func LoadAllOrders(db *sql.DB) ([]types.Order, error) {
	rows, err := db.Query("SELECT order_uid FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		o, err := LoadSingleOrder(db, id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
