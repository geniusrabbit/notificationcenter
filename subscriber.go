//
// @project geniusrabbit.com 2015 – 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016
//

package notificationcenter

// Subscriber data type
type Subscriber interface {
	// Subscribe new handler
	// @return error or nil
	Subscribe(h Handler) error

	// Unsubscribe this handler by ptr
	// @return error or nil
	Unsubscribe(h Handler) error

	// Start processing queue
	// @return error or nil
	Listen() error
}
