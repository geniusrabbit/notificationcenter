//
// @project geniusrabbit::notificationcenter 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

// Formater message converter
type Formater interface {
	Format(msg interface{}) interface{}
}

type FnkFormat func(msg interface{}) interface{}

func (f FnkFormat) Format(msg interface{}) interface{} {
	return f(msg)
}
