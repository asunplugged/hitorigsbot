package rimg

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/nfnt/resize"

	"git.trj.tw/golang/mtfosbot/module/config"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/utils"
)

func resizeImage(fio io.Reader, size int, wgth bool) (img image.Image, err error) {
	img, _, err = image.Decode(fio)
	w := 0
	h := 0
	if wgth {
		w = size
	} else {
		h = size
	}

	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return
}

// GetOriginImage -
func GetOriginImage(c *context.Context) {
	// max 1024
	fname := c.Param("imgname")
	conf := config.GetConf()
	imgP := utils.ParsePath(conf.ImageRoot)
	exists := utils.CheckExists(imgP, true)
	if !exists {
		c.NotFound("image path not found")
		return
	}

	subd := c.DefaultQuery("d", "")
	if len(subd) > 0 {
		imgP = path.Join(imgP, subd)
		exists = utils.CheckExists(imgP, true)
		if !exists {
			c.NotFound("image path not found")
			return
		}
	}

	newP := path.Join(imgP, fname)
	exists = utils.CheckExists(newP, false)
	if !exists {
		c.NotFound("image file not found")
		return
	}

	fileBuf, err := os.Open(newP)
	if err != nil {
		c.ServerError(nil)
		return
	}
	defer fileBuf.Close()

	imgf, format, err := image.DecodeConfig(fileBuf)
	if err != nil {
		c.ServerError(nil)
		return
	}

	// reset pointer to 0
	_, err = fileBuf.Seek(0, 0)
	if err != nil {
		c.ServerError(nil)
		return
	}

	if imgf.Height > 1024 || imgf.Width > 1024 {
		img, err := resizeImage(fileBuf, 1024, imgf.Width > imgf.Height)

		buf := new(bytes.Buffer)

		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			c.ServerError(nil)
			return
		}
		defer func() {
			buf.Reset()
			buf = nil
			img = nil
		}()

		c.Writer.Header().Set("Content-Type", "image/jpeg")
		c.Writer.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
		_, err = c.Writer.Write(buf.Bytes())
		if err != nil {
			c.ServerError(nil)
		}
		return
	}

	// check type
	if format != "jpeg" {
		img, _, err := image.Decode(fileBuf)
		if err != nil {
			c.ServerError(nil)
			return
		}
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			c.ServerError(nil)
			return
		}
		defer func() {
			buf = nil
			img = nil
		}()

		c.Writer.Header().Set("Content-Type", "image/jpeg")
		io.Copy(c.Writer, buf)
		return
	}

	c.Writer.Header().Set("Content-Type", "image/jpeg")
	io.Copy(c.Writer, fileBuf)
}

// GetThumbnailImage -
func GetThumbnailImage(c *context.Context) {
	// max 240
	fname := c.Param("imgname")
	conf := config.GetConf()
	imgP := utils.ParsePath(conf.ImageRoot)
	exists := utils.CheckExists(imgP, true)
	if !exists {
		c.NotFound("image path not found")
		return
	}

	thumbDir := "thumbnail"
	subd := c.DefaultQuery("d", "")
	if exists := utils.CheckExists(path.Join(imgP, thumbDir, subd), true); !exists {
		if err := os.MkdirAll(path.Join(imgP, thumbDir, subd), 0775); err != nil {
			c.ServerError(nil)
			return
		}
	}

	thumbP := path.Join(imgP, thumbDir, subd, fname)
	genNew := false
	if !utils.CheckExists(thumbP, false) {
		genNew = true
		thumbP = path.Join(imgP, subd, fname)
		exists = utils.CheckExists(thumbP, false)
		if !exists {
			c.NotFound("image file not found")
			return
		}
	}

	filebuf, err := os.Open(thumbP)
	if err != nil {
		c.ServerError(nil)
		return
	}
	defer filebuf.Close()

	imgconf, format, err := image.DecodeConfig(filebuf)
	if err != nil {
		c.ServerError(nil)
		return
	}

	// reset file reader ptr
	_, err = filebuf.Seek(0, 0)
	if err != nil {
		c.ServerError(nil)
		return
	}

	if imgconf.Height > 240 || imgconf.Width > 240 {
		img, err := resizeImage(filebuf, 240, imgconf.Width > imgconf.Height)
		if err != nil {
			c.ServerError(nil)
			return
		}

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			c.ServerError(nil)
			return
		}

		breader := bytes.NewReader(buf.Bytes())
		defer func() {
			buf = nil
			breader = nil
			img = nil
		}()

		if genNew {
			savep := path.Join(conf.ImageRoot, "thumbnail", subd, fname)
			err := saveNewThumbnail(breader, savep)
			if err != nil {
				c.ServerError(nil)
				return
			}
			_, err = breader.Seek(0, 0)
			if err != nil {
				c.ServerError(nil)
				return
			}
		}

		c.Writer.Header().Set("Content-Type", "image/jpeg")
		io.Copy(c.Writer, breader)
		return
	}

	if format != "jpeg" {
		img, _, err := image.Decode(filebuf)
		if err != nil {
			c.ServerError(nil)
			return
		}
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			c.ServerError(nil)
			return
		}

		breader := bytes.NewReader(buf.Bytes())
		defer func() {
			buf = nil
			breader = nil
			img = nil
		}()

		savep := path.Join(conf.ImageRoot, "thumbnail", fname)
		err = saveNewThumbnail(breader, savep)
		if err != nil {
			c.ServerError(nil)
			return
		}

		_, err = breader.Seek(0, 0)
		if err != nil {
			c.ServerError(nil)
			return
		}

		c.Writer.Header().Set("Content-Type", "image/jpeg")
		io.Copy(c.Writer, breader)
		return
	}

	if genNew {
		savep := path.Join(conf.ImageRoot, "thumbnail", fname)
		err := saveNewThumbnail(filebuf, savep)
		if err != nil {
			c.ServerError(nil)
			return
		}
		_, err = filebuf.Seek(0, 0)
		if err != nil {
			c.ServerError(nil)
			return
		}
	}

	c.Writer.Header().Set("Content-Type", "image/jpeg")
	io.Copy(c.Writer, filebuf)
}

func saveNewThumbnail(fio io.Reader, newp string) (err error) {
	out, err := os.Create(newp)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, fio)
	if err != nil {
		return err
	}

	return
}

// GetLineLogImage -
func GetLineLogImage(c *context.Context) {
	conf := config.GetConf()
	fname := c.Param("imgname")

	if !utils.CheckExists(path.Join(conf.LogImageRoot, fname), false) {
		c.NotFound("image not found")
		return
	}

	f, err := os.Open(path.Join(conf.LogImageRoot, fname))
	if err != nil {
		c.ServerError(nil)
		return
	}
	defer f.Close()

	r := io.LimitReader(f, 512)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		c.ServerError(nil)
		return
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		c.ServerError(nil)
		return
	}

	contentType := http.DetectContentType(bytes)

	c.Writer.Header().Set("Content-Type", contentType)
	io.Copy(c.Writer, f)
}
