package mu

func MSI(args ...interface{}) map[string]interface{} {
	output := make(map[string]interface{})

	if len(args)%2 != 0 {
		panic("invalid number of arguments")
	}

	for i := 0; i < len(args); i += 2 {
		output[args[i].(string)] = args[i+1]
	}

	return output
}
