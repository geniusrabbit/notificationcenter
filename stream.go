//
// @project geniusrabbit.com 2015 – 2016, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2019
//

package notificationcenter

// Streamer pipeline base declaration
type Streamer interface {
	// Send data to statistic
	Send(messages ...interface{}) error
}
