package model

type LineTwitchRT struct {
	Line   string `db:"line" cc:"line"`
	Twitch string `db:"twitch" cc:"twitch"`
	Tmpl   string `db:"tmpl" cc:"tmpl"`
	Type   string `db:"type" cc:"type"`
}

type LineFacebookRT struct {
	Line     string `db:"line" cc:"line"`
	Facebook string `db:"facebook" cc:"facebook"`
	Tmpl     string `db:"tmpl" cc:"tmpl"`
}

type LineYoutubeRT struct {
	Line    string `db:"line" cc:"line"`
	Youtube string `db:"youtube" cc:"youtube"`
	Tmpl    string `db:"tmpl" cc:"tmpl"`
}

type LineIGRT struct {
	Line string `db:"line" cc:"line"`
	IG   string `db:"ig" cc:"ig"`
	Tmpl string `db:"tmpl" cc:"tmpl"`
}

// AddRT - add facebook line rt
func (p *LineFacebookRT) AddRT() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_fb_rt" ("line", "facebook", "tmpl") values (:line, :facebook, :tmpl)`, p)
	return
}

// DelRT -
func (p *LineFacebookRT) DelRT() (err error) {
	_, err = x.NamedExec(`delete from "public"."line_fb_rt" where "line" = :line and "facebook" = :facebook`, p)
	return
}

// GetRT -
func (p *LineFacebookRT) GetRT() (err error) {
	stmt, err := x.PrepareNamed(`select * from "public"."line_fb_rt" where "line" = :line and "facebook" = :facebook`)
	if err != nil {
		return err
	}
	err = stmt.Get(p, p)
	return
}

// AddRT - add twitch line rt
func (p *LineTwitchRT) AddRT() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_twitch_rt" ("line", "twitch", "type", "tmpl") values (:line, :twitch, :type, :tmpl)`, p)
	return
}

// DelRT -
func (p *LineTwitchRT) DelRT() (err error) {
	_, err = x.NamedExec(`delete from "public"."line_twitch_rt" where "line" = :line and "twitch" = :twitch and "type" = :type`, p)
	return
}

// AddRT - add youtube line rt
func (p *LineYoutubeRT) AddRT() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_youtube_rt" ("line", "youtube", "tmpl") values (:line, :youtube, :tmpl)`, p)
	return
}

// DelRT -
func (p *LineYoutubeRT) DelRT() (err error) {
	_, err = x.NamedExec(`delete from "public"."line_youtube_rt" where "line" = :line and "youtube" = :youtube`, p)
	return
}

// AddRT -
func (p *LineIGRT) AddRT() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_ig_rt" ("line", "ig", "tmpl") values (:line, :ig, :tmpl)`, p)
	return
}

// DelRT -
func (p *LineIGRT) DelRT() (err error) {
	_, err = x.NamedExec(`delete from "public"."line_ig_rt" where "line" = :line and "ig" = :ig`, p)
	return
}
