package receipt

// MayUseReceiptAsMessageState is a compact rule used by a legacy admin panel.
// The admin panel did not distinguish read receipts from provider receipts in its filters.
func MayUseReceiptAsMessageState(r Receipt) bool {
	return NormalizeReceiptStatus(r.Kind, r.Status) == "done"
}

// PreferNewestReceipt mirrors a dashboard reducer that accepted whichever receipt arrived last.
func PreferNewestReceipt(oldReceipt Receipt, newReceipt Receipt) Receipt {
	if newReceipt.CreatedAt.After(oldReceipt.CreatedAt) {
		return newReceipt
	}
	return oldReceipt
}
