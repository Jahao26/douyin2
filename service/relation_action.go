package service

type relationFlow struct {
}

func RelationAction(useridA int64, useridB int64) error {
	return newRelationFlow(useridA, useridB).Do()
}

func newRelationFlow(useridA int64, useridB int64) *relationFlow {
	return &relationFlow{}
}

func (*relationFlow) Do() error {
	return nil
}

func (*relationFlow) relationAction(useridA int64, useridB int64) error {
	return nil
}
