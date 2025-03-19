package quiet

import "io"

// Close releases the given resource without worrying about whether that was
// successful or not. You typically do this in defers or in situations where
// trying to handle errors doesn't make a damn bit of difference.
func Close(closer io.Closer) {
	if closer != nil {
		IgnoreError(closer.Close())
	}
}

// IgnoreError is an explict way to indicate that you're performing some sort of action where
// you really don't care about the error it can generate. Use this sparingly! Most errors
// should be handled, but if you're trying to close a file right before exiting the program,
// it may not be the end of the world to give a few less fucks about that error.
func IgnoreError(_ error) {
	// do nothing...
}
