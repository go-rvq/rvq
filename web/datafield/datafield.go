package datafield

type Data map[any]any

func (d Data) GetData(key any) any {
	return d[key]
}

func (d *Data) SetData(key any, value any) {
	if d == nil {
		*d = Data{key: value}
	} else {
		(*d)[key] = value
	}
}

func (d Data) Clone() Data {
	if d == nil {
		return nil
	}

	newData := make(Data, len(d))
	for k, v := range d {
		newData[k] = v
	}

	return newData
}

type (
	DataFielder[T any] interface {
		GetDataField() *DataField[T]
	}

	DataField[T any] struct {
		dot  T
		data Data
	}
)

func New[T DataFielder[T]](dot T) T {
	dot.GetDataField().dot = dot
	return dot
}

func (d *DataField[T]) GetDataField() *DataField[T] {
	return d
}

func (d *DataField[T]) Data() Data {
	return d.data
}

func (d *DataField[T]) SetMapData(data Data) T {
	d.data = data
	return d.dot
}

func (d *DataField[T]) GetData(key any) any {
	return d.data[key]
}

func (d *DataField[T]) GetDataString(key any) (s string) {
	s, _ = d.GetData(key).(string)
	return
}

func (d *DataField[T]) GetDataInt(key any) (i int) {
	i, _ = d.GetData(key).(int)
	return
}

func (d *DataField[T]) GetDataBool(key any) (b bool) {
	b, _ = d.GetData(key).(bool)
	return
}

func (d *DataField[T]) SetData(key, value any) T {
	if d.data == nil {
		d.data = Data{}
	}
	d.data[key] = value
	return d.dot
}

func (d DataField[T]) Clone() DataField[T] {
	if d.data != nil {
		d.data = d.data.Clone()
	}
	return d
}
