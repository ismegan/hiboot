// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package websocket provides web socket auto configuration for web/cli application
package websocket

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/kataras/iris/websocket"
)

const (
	// Profile websocket profile name
	Profile = "websocket"
	// All is the string which the Emitter use to send a message to all.
	All = ""
	// Broadcast is the string which the Emitter use to send a message to all except this connection.
	Broadcast = ";to;all;except;me;"
)

type configuration struct {
	at.AutoConfiguration

	Properties properties `mapstructure:"websocket"`
}

// Connection is the websocket connection
type Connection interface {
	websocket.Connection
}

// ConnectionFunc is the websocket connection function
type ConnectionFunc func(ctx *web.Context, constructor HandlerConstructor) websocket.Connection

// ConnectionFunc is the websocket connection function
type HandlerConstructor func(conn Connection) Handler

type Server struct {
	*websocket.Server
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// Server websocket server
func (c *configuration) Server() *Server {
	s := websocket.New(websocket.Config{
		ReadBufferSize:  c.Properties.ReadBufferSize,
		WriteBufferSize: c.Properties.WriteBufferSize,
	})

	return &Server{
		Server: s,
	}
}

//func (c *configuration) Connection(server *Server, ctx *web.Context) websocket.Connection {
//	return server.Upgrade(ctx)
//}

// ConnectionFunc websocket connection for runtime dependency injection
func (c *configuration) ConnectionFunc(server *Server) ConnectionFunc {
	return func(ctx *web.Context, constructor HandlerConstructor) websocket.Connection {
		conn := server.Upgrade(ctx)
		handler := constructor(conn)
		conn.OnMessage(handler.OnMessage)
		conn.OnDisconnect(handler.OnDisconnect)
		conn.Wait()
		return conn
	}
}
