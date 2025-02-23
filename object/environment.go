package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

/*
Create a new environment.

Returns the new environment.
*/
func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{store: s}
}

/*
Create a new environment that is enclosed by the given environment.

Returns the new environment.
*/
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

/*
Get an object from the environment by name. If the object is not found in the current environment,
check the outer environment if it exists.

Returns the object, the environment we were searching in and a boolean indicating if the object was found.
*/
func (e *Environment) Get(name string) (Object, *Environment, bool) {
	env := e
	obj, ok := e.store[name]

	// If the object is not found in the current environment,
	// check the outer environment if it exists
	if !ok && e.outer != nil {
		obj, env, ok = e.outer.Get(name)
	}

	// If the object is found in the current environment or not found
	// return the object and the environment we were searching in
	return obj, env, ok
}

/*
Get an object from the current environment by name.

If the object is found in the current environment return it and true. Otherwise return the nil and false.
*/
func (e *Environment) GetFromCurrent(name string) (Object, bool) {
	obj, ok := e.store[name]

	return obj, ok
}

/*
Set the value of an identifier in the environment.

Returns the value that was set.
*/
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}
