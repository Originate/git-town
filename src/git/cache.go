package git

// BoolCache caches a string value.
// The zero value is an empty cache.
type BoolCache struct {
	value       bool
	initialized bool
}

// Set allows collaborators to signal when the current branch has changed.
func (sc *BoolCache) Set(newValue bool) {
	sc.value = newValue
	sc.initialized = true
}

// Value provides the current value.
func (sc *BoolCache) Value() bool {
	if !sc.initialized {
		panic("using current branch before initialization")
	}
	return sc.value
}

// Initialized indicates if we have a current branch.
func (sc *BoolCache) Initialized() bool {
	return sc.initialized
}

// Invalidate removes the cached value.
func (sc *BoolCache) Invalidate() {
	sc.initialized = false
}

// StringCache caches a string value.
// The zero value is an empty cache.
type StringCache struct {
	value       string
	initialized bool
}

// Set allows collaborators to signal when the current branch has changed.
func (sc *StringCache) Set(newValue string) {
	sc.value = newValue
	sc.initialized = true
}

// Value provides the current value.
func (sc *StringCache) Value() string {
	if !sc.initialized {
		panic("cannot access unitialized cached value")
	}
	return sc.value
}

// Initialized indicates if we have a current branch.
func (sc *StringCache) Initialized() bool {
	return sc.initialized
}

// Invalidate removes the cached value.
func (sc *StringCache) Invalidate() {
	sc.initialized = false
}

// StringSliceCache caches a string slice value.
// The zero value is an empty cache.
type StringSliceCache struct {
	value       []string
	initialized bool
}

// Set allows collaborators to signal when the current branch has changed.
func (ssc *StringSliceCache) Set(newValue []string) {
	ssc.value = newValue
	ssc.initialized = true
}

// Value provides the current value.
func (ssc *StringSliceCache) Value() []string {
	if !ssc.Initialized() {
		panic("cannot access unitialized cached value")
	}
	return ssc.value
}

// Initialized indicates if we have a current branch.
func (ssc *StringSliceCache) Initialized() bool {
	return ssc.initialized
}

// Invalidate removes the cached value.
func (ssc *StringSliceCache) Invalidate() {
	ssc.initialized = false
}
