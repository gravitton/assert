// Package assert provides a simple and lightweight testing assertion library for Go.
//
// This package offers a collection of assertion functions that can be used with any
// testing framework that implements the Testing interface (such as Go's standard
// testing.T). The assertions are designed to be intuitive and provide clear error
// messages when tests fail.
//
// # Basic Usage
//
// The assert package works with any type that implements the Testing interface:
//
//	func TestExample(t *testing.T) {
//		// Boolean assertions
//		assert.True(t, condition)
//		assert.False(t, !condition)
//
//		// Equality assertions
//		assert.Equal(t, actual, expected)
//		assert.NotEqual(t, actual, unexpected)
//
//		// Identity assertions (pointer comparison)
//		assert.Same(t, &actual, &expected)
//		assert.NotSame(t, &actual, &unexpected)
//	}
//
// # Assertion Functions
//
// • True/False: Assert boolean values
// • Equal/NotEqual: Assert value equality using reflect.DeepEqual
// • Same/NotSame: Assert pointer identity (same memory address)
// • Fail/Failf: Manually fail tests with custom messages
//
// # Error Messages
//
// All assertion functions provide detailed error messages when assertions fail,
// including actual and expected values formatted with Go's %#v verb for
// maximum clarity.
//
// # Helper Support
//
// All assertion functions automatically call t.helper() if the testing type
// implements the helper interface, ensuring that test failures point to the
// correct line in your test code rather than inside the assertion library.
//
// # Return Values
//
// All assertion functions return a boolean indicating success (true) or
// failure (false), allowing for conditional test logic if needed.
package assert
