package argparse

type ArgVariadic struct {
	Arg Parser
}

func (ArgVariadic) Empty() (interface{}, error) {
	return []interface{}{}, nil
}

func (arg ArgVariadic) Parse(str string) (interface{}, int, error) {
	values := []interface{}{}
	nparsed, inLen := 0, len(str)
	for {
		nparsed += countLeftSpaces(str[nparsed:])
		if nparsed >= inLen {
			break
		}
		val, n, err := arg.Arg.Parse(str[nparsed:])
		if err != nil {
			return nil, nparsed, err
		}
		nparsed += n
		values = append(values, val)
	}
	return values, nparsed, nil
}
