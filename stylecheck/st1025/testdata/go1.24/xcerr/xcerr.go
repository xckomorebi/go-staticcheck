package xcerr

import "fmt"

func New(obj string, meth string, format string) error {
	if obj == "" {
		return fmt.Errorf("%s: %s", meth, format)
	}
	return fmt.Errorf("%s.%s: %s", obj, meth, format)
}
