//
// @project geniusrabbit.com 2015, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015, 2019
//

package notificationcenter

import (
	"errors"
)

// Errors set
var (
	ErrInvalidObject              = errors.New("[notificationcenter] invalid handler")
	ErrInterfaceAlreadySubscribed = errors.New("[notificationcenter] interface already subscribed")
	ErrInvalidParams              = errors.New("[notificationcenter] invalid params")
)
