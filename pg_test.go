// Copyright (C) 2018 The Nori Authors info@nori.io
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program; if not, see <http://www.gnu.org/licenses/>.
package main

import (
	"testing"

	"github.com/cheebo/go-config"
	"github.com/go-pg/pg/orm"
	"github.com/stretchr/testify/assert"

	cfg "github.com/nori-io/nori/core/config"
)

const (
	testAddr     = "localhost:5432"
	testUser     = "postgres"
	testPassword = "postgres"
	testDB       = "postgres"
)

type TestPgTable struct {
	Id         int64
	TestField  string
	TestField2 string
}

func TestPlugin(t *testing.T) {
	assert := assert.New(t)

	cfgTest := go_config.New()
	cfgTest.SetDefault("pg.addr", testAddr)
	cfgTest.SetDefault("pg.user", testUser)
	cfgTest.SetDefault("pg.password", testPassword)
	cfgTest.SetDefault("pg.db", testDB)

	m := cfg.NewManager(cfgTest)

	p := new(plugin)

	assert.NotNil(p.Meta())
	assert.NotEmpty(p.Meta().Id())

	err := p.Init(nil, m)
	assert.Nil(err)

	err = p.Start(nil, nil)
	assert.Nil(err)

	db, ok := p.Instance().(Pg)
	assert.True(ok)
	assert.NotNil(db)

	table := new(TestPgTable)
	err = db.GetDB().CreateTable(table, &orm.CreateTableOptions{IfNotExists: true})
	assert.Nil(err)

	data := &TestPgTable{0, "1", "2"}
	err = db.GetDB().Insert(data)
	assert.Nil(err)

	data2 := new(TestPgTable)
	err = db.GetDB().Model(data2).Where("test_field = ?", "1").Select()
	assert.Nil(err)

	assert.Equal(data.TestField2, data2.TestField2)

	err = db.GetDB().DropTable(table, &orm.DropTableOptions{})
	assert.Nil(err)

	err = p.Stop(nil, nil)
	assert.Nil(err)
}
