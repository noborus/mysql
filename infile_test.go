// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2013 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"database/sql/driver"
	"testing"
	"time"
)

func Test_encodedLoadData(t *testing.T) {
	type args struct {
		x interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test String",
			args: args{"test"},
			want: "test",
		},
		{
			name: "test Int64",
			args: args{driver.Value(int64(42))},
			want: "42",
		},
		{
			name: "test Uint64",
			args: args{driver.Value(uint64(42))},
			want: "42",
		},
		{
			name: "test Fload64",
			args: args{driver.Value(float64(42.23))},
			want: "42.23",
		},
		{
			name: "test Bool",
			args: args{driver.Value(bool(true))},
			want: "1",
		},
		{
			name: "test BoolFalse",
			args: args{driver.Value(bool(false))},
			want: "0",
		},
		{
			name: "test nil",
			args: args{driver.Value(nil)},
			want: "\\N",
		},
		{
			name: "test TimeNULL",
			args: args{driver.Value(time.Time{})},
			want: "0000-00-00",
		},
		{
			name: "test Time",
			args: args{driver.Value(time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC))},
			want: "2014-12-31 12:13:24",
		},
		{
			name: "test byteNil",
			args: args{driver.Value([]byte(nil))},
			want: "\\N",
		},
		{
			name: "test byte",
			args: args{driver.Value([]byte("test"))},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &mysqlConn{
				cfg:              NewConfig(),
				maxAllowedPacket: defaultMaxAllowedPacket,
			}
			if got := mc.encodedLoadData(tt.args.x); got != tt.want {
				t.Errorf("mysqlConn.encodedLoadData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_escapedText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{"test"},
			want: "test",
		},
		{
			name: "test TAB",
			args: args{"t\test"},
			want: "t\\test",
		},
		{
			name: "test LF",
			args: args{"t\nest"},
			want: "t\\nest",
		},
		{
			name: "test CR",
			args: args{"t\rest"},
			want: "t\\rest",
		},
		{
			name: "test BackSlash",
			args: args{"t\\est"},
			want: "t\\\\est",
		},
		{
			name: "test All",
			args: args{"t\t\n\r\\est"},
			want: "t\\t\\n\\r\\\\est",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapedText(tt.args.text); got != tt.want {
				t.Errorf("escapedText() = %v, want %v", got, tt.want)
			}
		})
	}
}
