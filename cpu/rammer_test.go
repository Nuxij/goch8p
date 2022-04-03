package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRammer_AddDevice(t *testing.T) {
	ram := NewRAM(0x1000)
	rammer := NewRammer(0xF, []Device{ram})
	assert.Equal(t, rammer.devices[ram.UUID()], ram)
}

func TestRammer_SetRegion(t *testing.T) {
	ram := NewRAM(Size)
	type args struct {
		start        uint16
		size         uint16
		device       Device
		deviceOffset uint16
	}
	tests := []struct {
		name    string
		r       *Rammer
		args    args
		wantErr error
	}{
		{"first", NewRammer(0x100, []Device{ram}), args{0x0, 0x800, ram, 0x800}, nil},
		{"halfway", NewRammer(0x100, []Device{ram}), args{0x800, 0x200, ram, 0x0}, nil},
		{"bad alignment", NewRammer(0x100, []Device{ram}), args{0xAF, 0x800, ram, 0xF}, InvalidRegionAlignment{0xF, 0x0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.SetRegion(tt.args.start, tt.args.size, tt.args.device, tt.args.deviceOffset)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.size, tt.r.GetRegion(tt.args.start).Devices[0].Size)
				assert.Equal(t, ram.UUID(), tt.r.GetRegion(tt.args.start).Devices[0].ID)
				assert.Equal(t, tt.args.deviceOffset, tt.r.GetRegion(tt.args.start).Devices[0].Offset)
			}
		})
	}
}

func TestRammer_Read(t *testing.T) {
	type args struct {
		start        uint16
		size         uint16
		device       Device
		deviceOffset uint16
		readAddress  uint16
	}
	ram := NewRAM(Size)
	ram.Write(0x0, 0x90)
	ram.Write(0x800, 0xFF)
	ram.Write(0xEFF, 0x90)
	tests := []struct {
		name    string
		r       *Rammer
		args    args
		want    byte
		wantErr bool
	}{
		{"first", NewRammer(0x100, []Device{ram}), args{0x0, 0x200, ram, 0x0, 0x0}, 0x90, false},
		{"halfway", NewRammer(0x100, []Device{ram}), args{0x0, 0x200, ram, 0x700, 0x100}, 0xFF, false},
		{"mapped", NewRammer(0x100, []Device{ram}), args{0x200, 0x600, ram, 0x0, 0x7FF}, 0x0, false},
		{"unmapped", NewRammer(0x100, []Device{ram}), args{0x200, 0x600, ram, 0x0, 0x800}, 0x0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.SetRegion(tt.args.start, tt.args.size, tt.args.device, tt.args.deviceOffset)
			assert.NoError(t, err)
			got, err := tt.r.Read(tt.args.readAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rammer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rammer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRammer_Write(t *testing.T) {
	type args struct {
		addr  uint16
		value byte
	}
	ram := NewRAM(Size)
	tests := []struct {
		name    string
		r       *Rammer
		args    args
		wantErr error
	}{
		{"first", NewRammer(0x100, []Device{ram}), args{0x0, 0x90}, nil},
		{"halfway", NewRammer(0x100, []Device{ram}), args{0x400, 0xFF}, nil},
		{"mapped", NewRammer(0x100, []Device{ram}), args{0x200, 0x90}, nil},
		{"unmapped", NewRammer(0x100, []Device{ram}), args{0x800, 0x90}, AddressInvalid{0x800}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.SetRegion(0x0, 0x800, ram, 0x0)
			err := tt.r.Write(tt.args.addr, tt.args.value); 
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				got, err := tt.r.Read(tt.args.addr)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.value, got)
			}
			
		})
	}
}
