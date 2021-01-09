package background

import (
	"time"

	"/golang/mtfosbot/model"
	googleapis "/golang/mtfosbot/module/apis/google"
)

func checkYoutubeSubscribe() {
	e := time.Now().Unix() + (4 * 60 * 60)
	yt, err := model.GetYoutubeChannelsWithExpire(e)
	if err != nil || len(yt) == 0 {
		return
	}

	for _, v := range yt {
		googleapis.SubscribeYoutube(v.ID)
	}
}
