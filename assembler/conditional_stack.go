package assembler

// conditionalStack is shared by source preprocessing and assembly passes.
// Deferred frames keep both branches visible until layout can resolve them.
type conditionalStack struct {
	frames []conditionalFrame
}

type conditionalFrame struct {
	parentActive bool
	condition    bool
	elseSeen     bool
	deferred     bool
}

func (s *conditionalStack) Active() bool {
	active := true
	for _, frame := range s.frames {
		branchActive := true
		if !frame.deferred {
			branchActive = frame.condition
			if frame.elseSeen {
				branchActive = !frame.condition
			}
		}
		active = active && frame.parentActive && branchActive
	}
	return active
}

func (s *conditionalStack) Push(parentActive, condition, deferred bool) {
	s.frames = append(s.frames, conditionalFrame{
		parentActive: parentActive,
		condition:    condition,
		deferred:     deferred,
	})
}

func (s *conditionalStack) Current() *conditionalFrame {
	if len(s.frames) == 0 {
		return nil
	}
	return &s.frames[len(s.frames)-1]
}

func (s *conditionalStack) Pop() bool {
	if len(s.frames) == 0 {
		return false
	}
	s.frames = s.frames[:len(s.frames)-1]
	return true
}

func (s *conditionalStack) Complete() bool {
	return len(s.frames) == 0
}
