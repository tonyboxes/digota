//     Digota <http://digota.com> - eCommerce microservice
//     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published
//     by the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.
//
//     You should have received a copy of the GNU Affero General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/util"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"math/big"
)

const (
	WildcardScope Scope = "WILDCARD"
	PublicScope   Scope = "PUBLIC"
	WriteScope    Scope = "WRITE"
	ReadScope     Scope = "READ"
)

type (
	Client struct {
		Serial string
		Scopes []Scope
	}
	Role  string
	Scope string
)

type clientKey struct{}

var clients []Client

func New(c []config.Client) {
	for _, v := range c {
		var scopes []Scope
		for _, scope := range v.Scopes {
			scopes = append(scopes, Scope(scope))
		}
		clients = append(clients, Client{
			Serial: v.Serial,
			Scopes: scopes,
		})
	}
}

// NewContext store user in ctx and return new ctx.
func NewContext(ctx context.Context, serialId *big.Int) context.Context {
	var c *Client
	var err error
	if c, err = GetClient(util.BigIntToHex(serialId)); err != nil {
		return ctx
	}
	return context.WithValue(ctx, clientKey{}, c)
}

// FromContext returns the User stored in ctx, if any.
func FromContext(ctx context.Context) (*Client, bool) {
	u, ok := ctx.Value(clientKey{}).(*Client)
	return u, ok
}

func GetClient(serialId string) (*Client, error) {
	// search for user
	for _, c := range clients {
		if c.Serial == serialId {
			return &c, nil
		}
	}
	return nil, errors.New("Cant find client.")
}