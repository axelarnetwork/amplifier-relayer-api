package api

import (
	"fmt"
)

// enum for TokenManagerType as declared in solidity
// https://github.com/axelarnetwork/interchain-token-service/blob/cc6ef7282e18c6b2842cbf1098b06161e38ea32e/contracts/interfaces/ITokenManagerType.sol#L9-L17
const (
	solidityTokenManagerNativeInterchainToken uint8 = iota // 0
	solidityTokenManagerMintBurnFrom                       // 1
	solidityTokenManagerLockUnlock                         // 2
	solidityTokenManagerLockUnlockFee                      // 3
	solidityTokenManagerMintBurn                           // 4
)

// TokenManagerTypeFromSolidityEnum converts solidity TokenManagerType to API TokenManagerType
func TokenManagerTypeFromSolidityEnum(solidityTokenManagerType uint8) (TokenManagerType, error) {
	switch solidityTokenManagerType {
	case solidityTokenManagerNativeInterchainToken:
		return TokenManagerNativeInterchainToken, nil
	case solidityTokenManagerMintBurnFrom:
		return TokenManagerMintBurnFrom, nil
	case solidityTokenManagerLockUnlock:
		return TokenManagerLockUnlock, nil
	case solidityTokenManagerLockUnlockFee:
		return TokenManagerLockUnlockFee, nil
	case solidityTokenManagerMintBurn:
		return TokenManagerMintBurn, nil
	default:
		return *new(TokenManagerType), fmt.Errorf("invalid TokenManagerType: %d", solidityTokenManagerType)
	}
}
