package boot

func InitAll() (err error) {
	if err = InitConfig(); err != nil {
		return
	}

	InitLogger()

	if err = NewService(); err != nil {
		return err
	}

	return nil
}
