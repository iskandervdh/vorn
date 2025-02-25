package object

import "testing"

func TestEnvironment(t *testing.T) {
	env := NewEnvironment()

	expected := 5
	env.Set("foo", NewInteger(nil, int64(expected)))

	object, environment, exists := env.Get("foo")

	if !exists {
		t.Fatalf("env.Get did not return an object")
	}

	result, ok := object.(*Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", object, object)
	}

	if result.Value != int64(expected) {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	if environment != env {
		t.Errorf("environment is not the same as the one the object was set in")
	}

	object, exists = env.GetFromCurrent("foo")

	if !exists {
		t.Fatalf("env.GetFromCurrent did not return an object")
	}

	result, ok = object.(*Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", object, object)
	}

	if result.Value != int64(expected) {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	_, exists = env.GetFromCurrent("bar")

	if exists {
		t.Fatalf("env.GetFromCurrent returned an object")
	}

	_, _, exists = env.Get("bar")

	if exists {
		t.Fatalf("env.Get returned an object")
	}
}

func TestEnclosedEnvironment(t *testing.T) {
	env := NewEnvironment()

	expectedFoo := 5
	env.Set("foo", NewInteger(nil, int64(expectedFoo)))

	enclosed := NewEnclosedEnvironment(env)

	expectedBar := 10
	enclosed.Set("bar", NewInteger(nil, int64(expectedBar)))

	object, environment, exists := enclosed.Get("bar")

	if !exists {
		t.Fatalf("enclosed.Get did not return an object")
	}

	result, ok := object.(*Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", object, object)
	}

	if result.Value != int64(expectedBar) {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expectedBar)
	}

	if environment != enclosed {
		t.Errorf("environment is not the same as the one the object was set in")
	}

	object, exists = enclosed.GetFromCurrent("bar")

	if !exists {
		t.Fatalf("enclosed.GetFromCurrent did not return an object")
	}

	result, ok = object.(*Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", object, object)
	}

	if result.Value != int64(expectedBar) {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expectedBar)
	}

	_, exists = enclosed.GetFromCurrent("baz")

	if exists {
		t.Fatalf("enclosed.GetFromCurrent returned an object")
	}

	_, _, exists = enclosed.Get("baz")

	if exists {
		t.Fatalf("enclosed.Get returned an object")
	}

	object, environment, exists = enclosed.Get("foo")

	if !exists {
		t.Fatalf("enclosed.Get did not return an object")
	}

	result, ok = object.(*Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", object, object)
	}

	if result.Value != int64(expectedFoo) {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expectedFoo)
	}

	if environment != env {
		t.Errorf("environment is not the same as the one the object was set in")
	}
}
