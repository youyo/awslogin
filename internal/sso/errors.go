package sso

import (
	"errors"

	ssooidctypes "github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
)

// IsInvalidGrantError は err チェーンに InvalidGrantException が含まれるか判定する
func IsInvalidGrantError(err error) bool {
	var target *ssooidctypes.InvalidGrantException
	return errors.As(err, &target)
}
