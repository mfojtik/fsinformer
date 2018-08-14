package cache

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/mfojtik/fsinformer/pkg/types"
)

type testFile struct {
	name string
}

func (f *testFile) Name() string {
	return f.name
}

func (f *testFile) Stat() os.FileInfo {
	panic("implement me")
}

func (f *testFile) Content() []byte {
	panic("implement me")
}

func (f *testFile) ContentSum256() string {
	panic("implement me")
}

func mockFile() types.File {
	return &testFile{name: "/tmp/foo"}
}

func Test_syncMapStore_Add(t *testing.T) {
	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "default",
			fields:  fields{store: sync.Map{}},
			args:    args{obj: mockFile()},
			wantErr: false,
		},
		{
			name:    "not file",
			fields:  fields{store: sync.Map{}},
			args:    args{obj: "foo"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			if err := c.Add(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("syncMapStore.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_syncMapStore_Update(t *testing.T) {
	existingMap := sync.Map{}
	existingMap.Store("/tmp/foo", mockFile())

	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "default",
			fields:  fields{store: existingMap},
			args:    args{obj: mockFile()},
			wantErr: false,
		},
		{
			name:    "not exists",
			fields:  fields{store: sync.Map{}},
			args:    args{obj: mockFile()},
			wantErr: true,
		},
		{
			name:    "invalid",
			fields:  fields{store: sync.Map{}},
			args:    args{obj: "foo"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			if err := c.Update(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("syncMapStore.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_syncMapStore_Delete(t *testing.T) {
	existingMap := sync.Map{}
	existingMap.Store("/tmp/foo", mockFile())

	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "default",
			fields:  fields{store: existingMap},
			args:    args{obj: mockFile()},
			wantErr: false,
		},
		{
			name:    "not exists",
			fields:  fields{store: sync.Map{}},
			args:    args{obj: mockFile()},
			wantErr: true,
		},
		{
			name:    "invalid",
			fields:  fields{store: sync.Map{}},
			args:    args{obj: "foo"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			if err := c.Delete(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("syncMapStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_syncMapStore_List(t *testing.T) {
	existingMap := sync.Map{}
	existingMap.Store("/tmp/foo", mockFile())
	existingMap.Store("/tmp/bar", mockFile())

	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		{
			name:   "default",
			fields: fields{store: existingMap},
			want:   []interface{}{mockFile(), mockFile()},
		},
		{
			name:   "empty",
			fields: fields{store: sync.Map{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			if got := c.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("syncMapStore.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncMapStore_ListKeys(t *testing.T) {
	existingMap := sync.Map{}
	existingMap.Store("/tmp/foo", mockFile())

	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "default",
			fields: fields{store: existingMap},
			want:   []string{"/tmp/foo"},
		},
		{
			name:   "empty",
			fields: fields{store: sync.Map{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			if got := c.ListKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("syncMapStore.ListKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncMapStore_Get(t *testing.T) {
	existingMap := sync.Map{}
	item := mockFile()
	item2 := &testFile{name: "/tmp/bar"}

	existingMap.Store("/tmp/foo", item)

	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		want1   bool
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: existingMap,
			},
			args: args{
				obj: item,
			},
			want:    item,
			want1:   true,
			wantErr: false,
		},
		{
			name: "not exists",
			fields: fields{
				store: existingMap,
			},
			args: args{
				obj: item2,
			},
			want:    nil,
			want1:   false,
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				store: existingMap,
			},
			args: args{
				obj: "blah",
			},
			want:    nil,
			want1:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			got, got1, err := c.Get(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("syncMapStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("syncMapStore.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("syncMapStore.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_syncMapStore_GetByKey(t *testing.T) {
	existingMap := sync.Map{}
	existingMap.Store("/tmp/foo", mockFile())

	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		want1   bool
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: existingMap,
			},
			args: args{
				key: "/tmp/foo",
			},
			want:    mockFile(),
			want1:   true,
			wantErr: false,
		},
		{
			name: "not exists",
			fields: fields{
				store: existingMap,
			},
			args: args{
				key: "/tmp/bar",
			},
			want:    nil,
			want1:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			got, got1, err := c.GetByKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("syncMapStore.GetByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("syncMapStore.GetByKey() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("syncMapStore.GetByKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_syncMapStore_Replace(t *testing.T) {
	existingMap := sync.Map{}
	existingMap.Store("/tmp/foo", mockFile())
	item2 := &testFile{name: "bar"}

	type fields struct {
		mutex sync.Mutex
		store sync.Map
	}
	type args struct {
		items []interface{}
		in1   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: sync.Map{},
			},
			args: args{
				items: []interface{}{item2},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				store: sync.Map{},
			},
			args: args{
				items: []interface{}{"foo"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &syncMapStore{
				mutex: tt.fields.mutex,
				store: tt.fields.store,
			}
			if err := c.Replace(tt.args.items, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("syncMapStore.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
