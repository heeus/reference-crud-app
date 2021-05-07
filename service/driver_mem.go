/*
 * Copyright (c) 2019-present Heeus authors
 */

package service

import "fmt"

//MemoryDriver s.e.
type MemoryDriver struct {
	storage map[string]interface{}

	logger *Logger
}

//Name s.e.
func (d *MemoryDriver) Name() string {
	return "Memory drivwer"
}

//Info s.e.
func (d *MemoryDriver) Info() string {
	return "Light driver"
}

//Init s.e.
func (d *MemoryDriver) Init(args map[string]string) error {
	fmt.Println("memory driver initialized")

	d.storage = map[string]interface{}{}

	return nil
}

//Free s.e.
func (d *MemoryDriver) Free() error {
	fmt.Println("memory driver freed")
	return nil
}

//Read s.e.
func (d *MemoryDriver) Clean(r *DBRequest) *DBResponse {
	d.storage = map[string]interface{}{}

	return &DBResponse{Status: 200}
}

//Read s.e.
func (d *MemoryDriver) Read(r *DBRequest) *DBResponse {
	var records []*Record

	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewViews) > 0 {
		records = make([]*Record, len(r.ViewViews))

		for i, v := range r.ViewViews {
			rec, err := d.read(r.Partition, &v)

			if err != nil {
				return &DBResponse{Status: 400, Error: err.Error()}
			}

			if rec != nil {
				records[i] = rec
			}
		}
	}

	return &DBResponse{Status: 200, Records: records}
}

func (d *MemoryDriver) read(partition int64, view *ViewView) (*Record, error) {
	if partition < 0 {
		return nil, fmt.Errorf("record partiotion number malformed")
	}

	if view.ViewType == "" {
		return nil, fmt.Errorf("record ViewType malformed")
	}

	key, err := buildKey(view.PartitionKey, view.ClusterKey)

	if err != nil {
		return nil, err
	}

	if t := d.get(fmt.Sprintf("%v", partition), view.ViewType, key); t != nil {
		return &Record{
			Key:    key,
			Values: t.(map[string]interface{}),
		}, nil
	}

	return nil, nil
}

//Insert s.e.
func (d *MemoryDriver) Insert(r *DBRequest) *DBResponse {
	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewMods) > 0 {
		for _, v := range r.ViewMods {
			err := d.insert(r.Partition, &v)

			if err != nil {
				return &DBResponse{Status: 400, Error: err.Error()}
			}
		}
	}

	return &DBResponse{Status: 200}
}

func (d *MemoryDriver) insert(partition int64, view *ViewMod) error {
	if view.ViewType == "" {
		return fmt.Errorf("record ViewType name malformed")
	}

	key, err := buildKey(view.PartitionKey, view.ClusterKey)

	if err != nil {
		return err
	}

	pnum := fmt.Sprintf("%d", partition)

	r := d.get(pnum, view.ViewType, key)

	if r == nil {
		d.set(pnum, view.ViewType, key, view.Values)
	}

	return nil
}

//Update s.e.
func (d *MemoryDriver) Update(r *DBRequest) *DBResponse {
	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewMods) > 0 {
		for _, v := range r.ViewMods {
			err := d.update(r.Partition, &v)

			if err != nil {
				return &DBResponse{Status: 400, Error: err.Error()}
			}
		}
	}

	return &DBResponse{Status: 200}
}

func (d *MemoryDriver) update(partition int64, view *ViewMod) error {
	if partition < 0 {
		return fmt.Errorf("record partiotion number malformed")
	}

	if view.ViewType == "" {
		return fmt.Errorf("record table name malformed")
	}

	key, err := buildKey(view.PartitionKey, view.ClusterKey)

	if err != nil {
		return err
	}

	p := fmt.Sprintf("%v", partition)

	r := d.get(p, view.ViewType, key)

	if r == nil {
		return fmt.Errorf("Record with key %v not exists int partition %v table %v", key, partition, view.ViewType)
	}

	if len(view.Values) > 0 {
		newValues := r.(map[string]interface{})

		for k, v := range view.Values {
			newValues[k] = v
		}

		d.set(p, view.ViewType, key, newValues)
	}

	return nil
}

//Scan s.e.
func (d *MemoryDriver) Scan(r *DBRequest) *DBResponse {
	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	return &DBResponse{Status: 200, Error: "Scan method not implemented yet"}
}

//Delete s.e.
func (d *MemoryDriver) Delete(r *DBRequest) *DBResponse {

	if r == nil {
		return &DBResponse{Status: 400, Error: "wrong request data"}
	}

	if len(r.ViewViews) > 0 {

		for _, v := range r.ViewViews {
			err := d.delete(r.Partition, &v)

			if err != nil {
				return &DBResponse{Status: 400, Error: err.Error()}
			}
		}
	}

	return &DBResponse{Status: 200}
}

func (d *MemoryDriver) delete(partition int64, view *ViewView) error {
	if partition < 0 {
		return fmt.Errorf("record partiotion number malformed")
	}

	if view.ViewType == "" {
		return fmt.Errorf("record table name malformed")
	}

	key, err := buildKey(view.PartitionKey, view.ClusterKey)

	if err != nil {
		return err
	}

	p := fmt.Sprintf("%v", partition)
	r := d.get(p, view.ViewType, key)

	if r != nil {
		d.set(p, view.ViewType, key, nil)
	}

	return nil
}

func (d *MemoryDriver) get(partition string, table string, key string) interface{} {
	if p, ok := d.storage[partition]; ok {
		if t, ok := p.(map[string]interface{})[table]; ok {
			if v, ok := t.(map[string]interface{})[key]; ok {
				return v
			}
		}
	}

	return nil
}

/*
func (d *MemoryDriver) scan(partition int, table string, startKey string, count int) (map[string]Record, error) {
	ps := fmt.Sprintf("%v", partition)
	result := map[string]Record{}
	p, ok := d.storage[ps]
	search := false

	c := count

	if c <= 0 {
		return result, nil
	}

	if !ok {
		return nil, fmt.Errorf("no records in partition %v", partition)
	}

	t, ok := p.(map[string]interface{})[table]

	if !ok {
		return nil, fmt.Errorf("no records in partition %v table %v", partition, table)
	}

	for key, values := range t.(map[string]interface{}) {
		if key == startKey {
			search = true
		}

		if search {
			result[key] = Record{
				Key:    key,
				Values: values.(map[string]interface{}),
			}
		}

		c--

		if c == 0 {
			break
		}
	}

	return result, nil
}*/

func (d *MemoryDriver) set(partition string, table string, key string, values map[string]interface{}) {
	p, ok := d.storage[partition]

	if !ok {
		p = map[string]interface{}{}
		d.storage[partition] = p
	}

	t, ok := p.(map[string]interface{})[table]
	if !ok {
		t = map[string]interface{}{}
		p.(map[string]interface{})[table] = t
	}

	t.(map[string]interface{})[key] = values
}
