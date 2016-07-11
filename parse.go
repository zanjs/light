package main

func parse(file string) (m *Mapper, err error) {
	var docs map[string]string
	docs, err = parseDocs(file)
	if err != nil {
		return
	}

	_ = docs

	return
}

func parseDocs(file string) (docs map[string]string, err error) {

	return
}
