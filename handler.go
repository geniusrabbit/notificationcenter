//
// @project geniusrabbit.com 2015 – 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016
//

package notificationcenter

// Handler interface
type Handler interface {
	Handle(item interface{}) error
}

// FuncHandler type
type FuncHandler func(item interface{}) error

// Handle this item
func (f FuncHandler) Handle(item interface{}) error {
	return f(item)
}
