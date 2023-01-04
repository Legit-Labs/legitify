package parsing_utils

func ResolveAnnotation(customField interface{}) []string {
	retval := make([]string, 0)
	switch t := customField.(type) {
	case []interface{}:
		for _, enricherString := range t {
			switch ts := enricherString.(type) {
			case string:
				retval = append(retval, ts)
			}
		}
	case string:
		retval = append(retval, t)
	}

	return retval
}
