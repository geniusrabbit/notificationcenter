//
// @project geniusrabbit.com 2015 – 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016
//

package notificationcenter

// Logger base interface
type Logger interface {
	// Send data to statistic
	Send(messages ...interface{}) error
}
