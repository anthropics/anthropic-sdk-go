package apijson

// exactness ranks union variants by how well the input fit. A variant
// whose constants all match (constHit) is picked outright; among the
// rest, structural fit decides.
type exactness struct {
	constMatch    uint32
	constMismatch uint32

	requiredMissing uint32
	fieldErrors     uint32
	coercions       uint32
	extras          fitExtra
}

type fitExtra uint8

const (
	noExtras fitExtra = iota
	hasExtras
)

func (e *exactness) noteConstMatch()           { e.constMatch++ }
func (e *exactness) noteConstMismatch()        { e.constMismatch++ }
func (e *exactness) noteRequiredMissing(n int) { e.requiredMissing += uint32(n) }
func (e *exactness) noteFieldError()           { e.fieldErrors++ }
func (e *exactness) noteCoercion()             { e.coercions++ }
func (e *exactness) noteExtras()               { e.extras = hasExtras }

func (e exactness) constHit() bool {
	return e.constMismatch == 0 && e.constMatch > 0
}

func (e exactness) perfect() bool {
	e.constMatch = 0
	return e == exactness{}
}

func (e exactness) betterThan(o exactness) bool {
	if h, oh := e.constHit(), o.constHit(); h != oh {
		return h
	}
	if e.constMatch != o.constMatch {
		return e.constMatch > o.constMatch
	}
	if e.requiredMissing != o.requiredMissing {
		return e.requiredMissing < o.requiredMissing
	}
	if e.fieldErrors != o.fieldErrors {
		return e.fieldErrors < o.fieldErrors
	}
	if e.coercions != o.coercions {
		return e.coercions < o.coercions
	}
	return e.extras < o.extras
}

// absorb folds a child union's fit into its parent.
func (e *exactness) absorb(o exactness) {
	e.requiredMissing += o.requiredMissing
	e.fieldErrors += o.fieldErrors
	e.coercions += o.coercions
	if o.extras > e.extras {
		e.extras = o.extras
	}
}
