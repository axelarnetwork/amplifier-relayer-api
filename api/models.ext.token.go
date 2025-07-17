package api

// GetTokenID returns the TokenID of the token.
func (t Token) GetTokenID() *string {
	return t.TokenID
}

// GetAmount returns the Amount of the token.
func (t Token) GetAmount() string {
	return t.Amount
}

// GetTokenID returns the TokenID of the token.
func (t UnsignedToken) GetTokenID() *string {
	return t.TokenID
}

// GetAmount returns the Amount of the token.
func (t UnsignedToken) GetAmount() string {
	return t.Amount
}
