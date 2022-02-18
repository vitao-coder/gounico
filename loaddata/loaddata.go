package loaddata

type LoadData interface {
	ProcessCSVToDatabase(csvByteArray []byte) error
}
