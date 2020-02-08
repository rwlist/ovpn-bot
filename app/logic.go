package app

type Logic struct {

}

func NewLogic() *Logic {
	return &Logic{}
}

func (l *Logic) CommandInit() (string, error) {
	return "ok", nil
}

func (l *Logic) CommandStatus() (string, error) {
	return "ok", nil
}

func (l *Logic) CommandGenerate(profileName string) (string, error) {
	return "ok", nil
}
