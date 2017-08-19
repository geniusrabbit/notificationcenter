//
// @project geniusrabbit::notificationcenter 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package metrics

import "strings"

func influxFormatFnk(msg interface{}) interface{} {
	switch m := msg.(type) {
	case Message:
		msg = influxFormatMsg(&m)
	case *Message:
		msg = influxFormatMsg(m)
	}
	return msg
}

func influxFormatMsg(msg *Message) interface{} {
	var key = msg.Name
	{
		var tags []string
		for k, v := range msg.Tags {
			tags = append(tags, k+`=`+v)
		}
		if len(tags) > 0 {
			key = msg.Name + `,` + strings.Join(tags, `,`)
		}
	}

	switch msg.Type {
	case MessageTypeIncrement:
		return key
	case MessageTypeGauge:
		return TypeGauge{key: msg.Value}
	case MessageTypeTiming:
		return TypeGauge{key: msg.ValueInt()}
	case MessageTypeCount:
		return TypeGauge{key: msg.ValueInt()}
	case MessageTypeUnique:
		return TypeGauge{key: msg.ValueString()}
	}
	return nil
}

// InfluxFormat message converter
var InfluxFormat = FnkFormat(influxFormatFnk)
