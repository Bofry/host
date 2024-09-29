package sdk

func Ensure(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func Assert(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
