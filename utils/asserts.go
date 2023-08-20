package utils

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MustExist(exists bool) {
	if !exists {
		panic("doesn't exists")
	}
}

func CheckErrMustExist(err error, exists bool) {
	CheckErr(err)
	MustExist(exists)
}
