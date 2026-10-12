package engine

// Stream-M05 Error Tags & Misconceptions for Module 5: Inventory & Cost of Goods Sold.
const (
	TagM05LifoUsedInIFRS              = "m05_lifo_used_in_ifrs"
	TagM05FreightInExpensed           = "m05_freight_in_expensed"
	TagM05NRVWriteDownReversed        = "m05_nrv_writedown_reversed"
	TagM05InventoryRecordedAtSelling  = "m05_inventory_recorded_at_selling_price"
	TagM05COGSConfusedWithOperating   = "m05_cogs_confused_with_operating_expense"
	TagM05InventoryCreditedOnPurchase = "m05_inventory_credited_on_purchase"
)

// Stream-M05 Family ID constants for candidate inventory questions.
const (
	FamilyM05InventoryFIFO         = "m05_inv_fifo_costing"
	FamilyM05InventoryLIFO         = "m05_inv_lifo_costing"
	FamilyM05InventoryNRVWriteDown = "m05_inv_nrv_writedown"
	FamilyM05FreightInCapitalize   = "m05_inv_freight_in_cap"
)
