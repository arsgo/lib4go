package influxdb

type configHandler interface {
	GetSourceConfig(string, string) (string, error)
}

func SaveMapsToInfluxDB(typeName string, fname string, rows []map[string]interface{}, handler configHandler) (err error) {
	config, err := handler.GetSourceConfig(typeName, fname)
	if err != nil {
		return
	}
	return Save(config, rows)
}
