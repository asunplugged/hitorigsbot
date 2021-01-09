package background

import (
	"github.com/robfig/cron"
)

var c *cron.Cron

// SetBackground -
func SetBackground() {
	c = cron.New()
	// c.AddFunc("0 * * * * *", readFacebookPage)
	c.AddFunc("*/20 * * * * *", getStreamStatus)
	c.AddFunc("*/5 * * * * *", checkOpay)
	c.AddFunc("0 0 */3 * * *", checkYoutubeSubscribe)
	c.Start()
}
