package util

func Assert(condition bool, msg string) {
	if !condition {
		panic(msg)
	}
}

func AssertNoErr(err error) {
	if err != nil {
		panic(err)
	}
}
