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
// along with this program; if not, see <http://www.gnu.org/licenses/>
package main

import (
	"context"

	"github.com/go-pg/pg"

	cfg "github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori-common/meta"
	noriPlugin "github.com/nori-io/nori-common/plugin"
)

type Pg interface {
	GetDB() *pg.DB
}

type plugin struct {
	instance Pg
	config   *config
}

type config struct {
	address  cfg.String
	user     cfg.String
	password cfg.String
	db       cfg.String
}

type instance struct {
	database *pg.DB
}

var (
	Plugin plugin
)

func (p *plugin) Init(_ context.Context, configManager cfg.Manager) error {
	m := configManager.Register(p.Meta())
	p.config = &config{
		address:  m.String("pg.addr", ""),
		user:     m.String("pg.user", ""),
		password: m.String("pg.password", ""),
		db:       m.String("pg.db", ""),
	}
	return nil
}

func (p *plugin) Instance() interface{} {
	return p.instance
}

func (p plugin) Meta() meta.Meta {
	return &meta.Data{

		ID: meta.ID{
			ID:      "nori/db/postgresql/pg",
			Version: "1.0.0",
		},
		Author: meta.Author{
			Name: "Nori",
			URI:  "https://nori.io",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name: "Nori: PostgreSQL ORM Driver",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"db", "postgresql", "orm"},
	}

}

func (p *plugin) Start(ctx context.Context, _ noriPlugin.Registry) error {
	if p.instance == nil {
		instance := &instance{
			database: pg.Connect(&pg.Options{
				Addr:     p.config.address(),
				User:     p.config.user(),
				Password: p.config.password(),
				Database: p.config.db(),
			}),
		}

		p.instance = instance
	}
	return nil
}

func (p *plugin) Stop(_ context.Context, _ noriPlugin.Registry) error {
	err := p.instance.(*instance).database.Close()
	p.instance = nil
	return err
}

func (i instance) GetDB() *pg.DB {
	return i.database
}
