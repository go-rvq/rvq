package presets

type ModelWalkState uint8

const (
	ModelWalkStateNext ModelWalkState = iota
	ModelWalkStateStop
	ModelWalkStateSkipChildren
)

func WalkModels(models []*ModelBuilder, f func(mb *ModelBuilder) (state ModelWalkState, err error)) (err error) {
	_, err = modelWalk(models, f)
	return
}

func modelWalk(models []*ModelBuilder, f func(mb *ModelBuilder) (state ModelWalkState, err error)) (state ModelWalkState, err error) {
	for _, mb := range models {
		if state, err = f(mb); err != nil {
			return
		}
		switch state {
		case ModelWalkStateNext:
			var state2 ModelWalkState
			if state2, err = modelWalk(mb.children, f); err != nil {
				return
			}
			switch state2 {
			case ModelWalkStateStop:
				state = ModelWalkStateStop
				return
			}
		case ModelWalkStateStop:
			return
		case ModelWalkStateSkipChildren:
			state = ModelWalkStateNext
			return
		}
	}
	return
}
