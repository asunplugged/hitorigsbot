package apimsg

// ResObject -
type ResObject struct {
	Status int
	Obj    interface{}
}

var objs = map[string]*ResObject{
	"NotFound": &ResObject{
		Status: 404,
		Obj: map[string]string{
			"message": "not found",
		},
	},
	"InternalError": &ResObject{
		Status: 500,
		Obj: map[string]string{
			"message": "server internal error",
		},
	},
	"Success": &ResObject{
		Status: 200,
		Obj: map[string]string{
			"message": "success",
		},
	},
	"Forbidden": &ResObject{
		Status: 403,
		Obj: map[string]string{
			"message": "forbidden",
		},
	},
	"DataFormat": &ResObject{
		Status: 400,
		Obj: map[string]string{
			"message": "input data format error",
		},
	},
	"LoginFirst": &ResObject{
		Status: 401,
		Obj: map[string]string{
			"message": "login first",
		},
	},
}

// GetRes -
func GetRes(name string, msg interface{}) *ResObject {
	obj, ok := objs[name]
	if !ok {
		obj = objs["InternalError"]
	}

	resobj := &ResObject{}
	resobj.Status = obj.Status
	switch msg.(type) {
	case string:
		tmp := make(map[string]string)
		tmp["message"] = msg.(string)
		resobj.Obj = tmp
		break
	case map[string]interface{}:
		resobj.Obj = msg
		break
	case map[string]string:
		resobj.Obj = msg
		break
	default:
		resobj.Obj = obj.Obj
	}
	return resobj
}
