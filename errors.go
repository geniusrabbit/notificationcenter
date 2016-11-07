//
// @project geniusrabbit.com 2015
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015
//

package notificationcenter

import (
	"errors"
)

// Errors set
var (
	ErrInvalidObject              = errors.New("Invalid handler")
	ErrInterfaceAlreadySubscribed = errors.New("Interface already subscribed")
	ErrInvalidParams              = errors.New("Invalid params")
)
