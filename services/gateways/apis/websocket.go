package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bridgekit-io/frodo/fail"
	"github.com/bridgekit-io/frodo/internal/quiet"
	"github.com/bridgekit-io/frodo/internal/radix"
	"github.com/bridgekit-io/frodo/services"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// WalkWebsockets invokes your callback/handler on all registered websockets that match the given connection
// ID prefix. For instance, if you want to push a notification to all of a user's active browser sessions,
// you could walk w/ the connection prefix "user.123" and that would push the notification to the connections
// named "user.123.A", "user.123.B", and "user.123.C".
//
// All callbacks are executed in separate goroutines, so expect these to run in parallel for all matching sockets.
func WalkWebsockets(ctx context.Context, socketPrefix string, handler func(ctx context.Context, websocket *Websocket) error) error {
	registry, ok := ctx.Value(websocketRegistryContextKey{}).(*websocketRegistry)
	if !ok {
		return fail.Unexpected("error connecting websocket: missing websocket registry")
	}

	if socketPrefix == "" {
		return fail.BadRequest("walking websockets requires a non-empty prefix")
	}

	errs, _ := fail.NewGroup(ctx)
	registry.walk(socketPrefix, func(websocket *Websocket) {
		errs.Go(func() error {
			// Ignore errors where the socket was already closed. Yes, the startListening() loop should auto-close the socket once the
			// conn is bad, but it's possible that this handler fires before it can fully break out of the loop, so let's give it a hand.
			switch err := handler(ctx, websocket); {
			case errors.Is(err, net.ErrClosed):
				quiet.Close(websocket)
				return nil
			default:
				return err
			}
		})
	})
	return errs.Wait()

}

// ConnectWebsocket hijacks the HTTP connection and makes it so that the user can have duplex communication with
// a connected client browser/device.
func ConnectWebsocket(ctx context.Context, connectionID string, opts WebsocketOptions) (*Websocket, error) {
	sockets, ok := ctx.Value(websocketRegistryContextKey{}).(*websocketRegistry)
	if !ok {
		return nil, fail.Unexpected("error connecting websocket: missing websocket registry")
	}
	req, ok := ctx.Value(requestContextKey{}).(*http.Request)
	if !ok {
		return nil, fail.Unexpected("error connecting websocket: missing request")
	}
	w, ok := ctx.Value(responseContextKey{}).(http.ResponseWriter)
	if !ok {
		return nil, fail.Unexpected("error connecting websocket: missing response")
	}

	// What the hell are we doing here? In Javascript, you can't customize headers on the ws/wss request when trying to
	// open a websocket. Unfortunately, that means we have to use "less pretty" ways to give users a way to do this.
	ws.DefaultHTTPUpgrader.Protocol = func(value string) bool {
		return strings.HasPrefix(value, "Authorization.")
	}

	// Upgrade the HTTP connection to a websocket.
	conn, _, _, err := ws.UpgradeHTTP(req, w)
	if err != nil {
		return nil, fmt.Errorf("error connecting websocket: %w", err)
	}

	// Message handlers should be able to walk the websocket registry.
	newMessageContext := func() context.Context {
		return context.WithValue(context.Background(), websocketRegistryContextKey{}, sockets)
	}

	// Make sure that we automatically clean up the registry when connections are lost or explicitly closed.
	opts = opts.assignDefaults(func() {
		sockets.remove(connectionID)
	})
	websocket := Websocket{Conn: conn, ID: connectionID, Options: opts, newMessageContext: newMessageContext}

	// Make sure that we only have one connection for a given connection id. If the same user can be connected in
	// different places, don't use the connection ID "user.123.web" because if they're logged in w/ 2 different
	// browsers they'll clobber each other. Instead, do "user.123.web.TIMESTAMP" or something like that. They'll
	// each maintain connections, and you can locate them both by doing prefix lookups on "user.123.web"
	if old, ok := sockets.add(connectionID, &websocket); ok {
		quiet.Close(old)
	}

	websocket.startListening()
	return &websocket, nil
}

// WebsocketOptions provides the necessary callbacks for handling incoming data read from client connections.
type WebsocketOptions struct {
	// Logger lets you customize how you want this debug/trace logging to work.
	Logger *slog.Logger
	// OnReadText provides a handler for incoming text data.
	OnReadText func(ctx context.Context, socket *Websocket, data []byte)
	// OnReadBinary provides a handler for incoming binary data.
	OnReadBinary func(ctx context.Context, socket *Websocket, data []byte)
	// OnReadContinuation provides a handler for incoming continuation frames.
	OnReadContinuation func(ctx context.Context, socket *Websocket, data []byte)
	// OnClose provides a custom handler that fires when this websocket is closed for any reason. This includes you manually calling
	// the Close() method on the Websocket or the server automatically closing a seemingly dead connection.
	OnClose func(ctx context.Context, socketID string)

	// PingInterval determines how frequently the server will send op-ping frames to the client in hopes of receiving
	// an op-pong in response. This ping/pong is how the server detects disconnected clients that abruptly ended their
	// socket connection (e.g. lost network or lost power) without gracefully sending an op-close to let the server know
	// about the intention to disconnect.
	//
	// This value allows you to control the tradeoff of network chatter vs realtime awareness. For instance, it may be okay
	// to only send a ping every 30 seconds. True, you might go along for 29 seconds thinking a connection is fine and dandy
	// when in reality the client dropped off, but this lets you prioritize chatter vs realtime.
	//
	// The default interval is 10s.
	PingInterval time.Duration
	// IdleTimeout indicates how long the server will keep a websocket around without receiving some sort of message or control
	// frame from the client. The socket's lease can be extended by receiving a successful "pong" in response to a recent
	// ping (see PingInterval for more details). Additionally, receiving a standard text/binary message will reset this since it's
	// proof that the underlying TCP connection is still alive and well.
	//
	// If the server does not receive a message or a pong or any other type of data frame from the client after this amount of time,
	// the server will close the underlying TCP connection, release any underlying resources, and invoke the OnClose handler you
	// configure in these options.
	//
	// It is a good idea to set this to at least 2 to 3 times as long as PingInterval. This way, if a ping frame is somehow lost,
	// you have more opportunities to get a sign of life from the client connection before killing the socket.
	IdleTimeout time.Duration
}

func (opts WebsocketOptions) assignDefaults(frameworkCleanup func()) WebsocketOptions {
	if opts.PingInterval <= 0 {
		opts.PingInterval = 5 * time.Second
	}
	if opts.IdleTimeout <= 0 {
		opts.IdleTimeout = 15 * time.Second
	}
	if opts.OnReadText == nil {
		opts.OnReadText = func(ctx context.Context, socket *Websocket, data []byte) {}
	}
	if opts.OnReadBinary == nil {
		opts.OnReadBinary = func(ctx context.Context, socket *Websocket, data []byte) {}
	}
	if opts.OnReadContinuation == nil {
		opts.OnReadContinuation = func(ctx context.Context, socket *Websocket, data []byte) {}
	}

	// We want OnClose to actually run both, our internal registry cleanup handler AND whatever logic the
	// user provides. Replace OnClose with a function that executes them both.
	switch userCleanup := opts.OnClose; userCleanup {
	case nil:
		opts.OnClose = func(_ context.Context, _ string) {
			frameworkCleanup()
		}
	default:
		opts.OnClose = func(ctx context.Context, socketID string) {
			frameworkCleanup()
			userCleanup(ctx, socketID)
		}
	}
	return opts
}

type websocketState int

// We only create Websocket structs AFTER successfully upgrading the connection, so its zero starting value should be "connected".
const (
	websocketStateConnected websocketState = iota
	websocketStateClosing
	websocketStateClosed
)

// Websocket wraps the connection and other necessary information to provide everything you need to communicate
// with the client over a recently opened websocket.
type Websocket struct {
	// ID is the unique identifier for this socket. Frodo provides a way for you to look up individual connections,
	// so you can push messages as needed, so this id/key helps you locate the websocket later.
	ID string
	// Conn is the actual TCP connection we're keeping open to handle the communication.
	Conn net.Conn
	// Options contains our callbacks for handling all manner of reads and the close event.
	Options WebsocketOptions
	// newMessageContext is used internally to create a context intended to be used for the handling of a single message written to the socket.
	newMessageContext func() context.Context
	// state helps determine where in the very simple state machine this socket resides (e.g. connected/closed).
	state websocketState
}

// WriteText writes a frame of binary data to the client on the other end of the socket.
func (socket *Websocket) Write(data []byte) error {
	if err := socket.assertConnected(); err != nil {
		return err
	}

	if err := wsutil.WriteServerBinary(socket.Conn, data); err != nil {
		quiet.Close(socket)
		return fmt.Errorf("error writing to websocket: %s: %w", socket.ID, err)
	}
	return nil
}

// WriteClose pushes a "close" frame to the client, letting them know that we want to close up shop.
func (socket *Websocket) WriteClose(data []byte) error {
	if err := socket.assertConnected(); err != nil {
		return err
	}

	if err := wsutil.WriteServerMessage(socket.Conn, ws.OpClose, data); err != nil {
		quiet.Close(socket)
		return fmt.Errorf("error writing to websocket: %s: %w", socket.ID, err)
	}
	return nil
}

// WriteText writes a frame of text data to the client on the other end of the socket.
func (socket *Websocket) WriteText(data string) error {
	if err := socket.assertConnected(); err != nil {
		return err
	}

	if err := wsutil.WriteServerText(socket.Conn, []byte(data)); err != nil {
		quiet.Close(socket)
		return fmt.Errorf("error writing to websocket: %s: %w", socket.ID, err)
	}
	return nil
}

// WriteJSON marshals the given object into a JSON string and then writes a text frame to the client.
func (socket *Websocket) WriteJSON(value any) error {
	if err := socket.assertConnected(); err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error writing to websocket: %w", err)
	}
	if err := wsutil.WriteServerText(socket.Conn, data); err != nil {
		quiet.Close(socket)
		return fmt.Errorf("error writing to websocket: %s: %w", socket.ID, err)
	}
	return nil
}

// assertConnected returns an appropriately phrased error message if the state of the socket is anything other than "connected".
func (socket *Websocket) assertConnected() error {
	switch socket.state {
	case websocketStateClosing:
		return fail.Unavailable("socket closing")
	case websocketStateClosed:
		return fail.Unavailable("socket closed")
	default:
		return nil
	}
}

// Close kills the current connection. This will also trigger your OnClose handler.
func (socket *Websocket) Close() error {
	if socket.state != websocketStateConnected { // make sure another goroutine isn't already trying to close this.
		return nil
	}

	socket.state = websocketStateClosing
	defer func() { socket.state = websocketStateClosed }()

	if socket.Options.Logger != nil {
		socket.Options.Logger.Info("Websocket CLOSED", "websocket_id", socket.ID)
	}

	// The context.Background() is to provide the callback with a context that is not bound to any now-closed HTTP/TCP
	// resources. It's a context dedicated to the user's callback and the task of cleanup.
	quiet.Close(socket.Conn)
	socket.Conn = nil
	socket.Options.OnClose(context.Background(), socket.ID)
	return nil
}

// startListening fires off a separate goroutine that infinitely loops, attempting to read messages from the client.
// Any incoming messages will be routed to the socket's OnReadText/Binary/etc. handlers. This will exit automatically
// when the socket/connection is closed.
func (socket *Websocket) startListening() {
	if socket.Options.Logger != nil {
		socket.Options.Logger.Info("Websocket OPENED", "websocket_id", socket.ID)
	}

	// Goroutine 1: Send OpPing frames at the configured interval and trigger the auto-close if we detect a dead connection.
	go func() {
		for socket.state == websocketStateConnected {
			err := wsutil.WriteServerMessage(socket.Conn, ws.OpPing, ws.CompiledPing)
			switch {
			case err != nil:
				quiet.Close(socket)
			default:
				time.Sleep(socket.Options.PingInterval)
			}
		}
	}()

	// Goroutine 2: Iterate and read all incoming frames/messages sent by the client. Each iteration is bound by the IdleTimeout to
	// detect a seemingly-dead TCP connection.
	go func() {
		defer quiet.Close(socket)

		var data []byte
		var frameReader = wsutil.Reader{
			Source:          socket.Conn,
			State:           ws.StateServerSide,
			CheckUTF8:       true,
			SkipHeaderCheck: false,
			OnIntermediate:  wsutil.ControlFrameHandler(socket.Conn, ws.StateServerSide),
		}

		for socket.state == websocketStateConnected {
			now := time.Now()
			deadline := now.Add(socket.Options.IdleTimeout)

			// The user configures how long they're willing to let a connection linger with no communication from the client via
			// the IdleTimeout. If we don't get a message or a pong or anything from them, we're going to close the connection
			// because the client probably disconnected in a hurry (lost internet, powered off, etc.).
			//
			// While we get to ignore a lot of failures in this block, don't ignore this one! This deadline is the only way we
			// don't let this connection just block forever with no chance of ever being used again - i.e. it will leak if there
			// is no limit to how long we'll try. Trust me, this will work. It's a *net.TCPConn, and it supports setting deadlines,
			// so as long as the connection is valid it will work. If this fails, the underlying TCP connection probably isn't
			// usable anyway, and we just need to have the client reconnect.
			if err := socket.Conn.SetDeadline(deadline); err != nil {
				quiet.Close(socket)
				continue
			}

			header, err := frameReader.NextFrame()
			switch {
			case errors.Is(err, os.ErrDeadlineExceeded):
				// The client isn't sending messages or responding to pings, and it exceeded our IdleTimeout. Assume the connection is lost.
				quiet.Close(socket)

			case err != nil:
				// Chances are the underlying *net.TCPConn is closed or hosed. Either way, the best solution is to let the client reconnect.
				quiet.Close(socket)

			/****** Control Op Codes ******/

			case header.OpCode == ws.OpClose:
				// The client either gets our ACK or it doesn't. Who cares - we're closing the socket anyway.
				quiet.IgnoreError(frameReader.OnIntermediate(header, &frameReader))
				quiet.Close(socket)

			case header.OpCode == ws.OpPing:
				// Try to send an OpPong in response, but don't stress over it. If we lost the connection, we'll know soon enough.
				quiet.IgnoreError(frameReader.OnIntermediate(header, &frameReader))

			case header.OpCode == ws.OpPong:
				// We're just discarding the response, so don't worry about failure. We got some sign of life from the
				// client, so we'll keep the connection alive... for now.
				quiet.IgnoreError(frameReader.OnIntermediate(header, &frameReader))

			/****** Message Op Codes ******/

			case header.OpCode == ws.OpText:
				if data, err = io.ReadAll(&frameReader); err != nil {
					quiet.Close(socket)
					continue
				}
				socket.Options.OnReadText(socket.newMessageContext(), socket, data)

			case header.OpCode == ws.OpBinary:
				if data, err = io.ReadAll(&frameReader); err != nil {
					quiet.Close(socket)
					continue
				}
				socket.Options.OnReadBinary(socket.newMessageContext(), socket, data)

			case header.OpCode == ws.OpContinuation:
				if data, err = io.ReadAll(&frameReader); err != nil {
					quiet.Close(socket)
					continue
				}
				socket.Options.OnReadContinuation(socket.newMessageContext(), socket, data)

			/****** Da-Fuck-You-Doing? Op Codes ******/

			default:
				socket.Options.Logger.Warn(fmt.Sprintf("Unexpected websocket op code '%x'", header.OpCode))
				quiet.IgnoreError(frameReader.Discard())
			}
		}
	}()
}

type websocketRegistryContextKey struct{}

// websocketRegistryMiddleware ensures that WalkWebsockets and ConnectWebsocket have access to the gateway's
// master websocket connection registry. We don't expose this registry to end users of Frodo. We only provide
// functions that let them indirectly interact w/ it.
func websocketRegistryMiddleware(websockets *websocketRegistry) services.MiddlewareFunc {
	return func(ctx context.Context, req any, next services.HandlerFunc) (any, error) {
		ctx = context.WithValue(ctx, websocketRegistryContextKey{}, websockets)
		return next(ctx, req)
	}
}

func newWebsocketRegistry() *websocketRegistry {
	return &websocketRegistry{
		mutex:   &sync.Mutex{},
		sockets: radix.New[*Websocket](),
	}
}

// websocketRegistry contains a mapping of all currently-open connections and their IDs.
type websocketRegistry struct {
	mutex   *sync.Mutex
	sockets radix.Tree[*Websocket]
}

func (registry *websocketRegistry) add(socketID string, socket *Websocket) (*Websocket, bool) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	return registry.sockets.Insert(socketID, socket)
}

func (registry *websocketRegistry) remove(socketID string) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	registry.sockets.Delete(socketID)
}

// walk invokes the visitorFunc on all websockets whose ids start w/ the given prefix. Beware! This locks
// the entire registry while doing this, so make sure that your visitor is not doing anything intensive.
func (registry *websocketRegistry) walk(socketPrefix string, visitorFunc func(*Websocket)) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	registry.sockets.WalkPrefix(socketPrefix, func(_ string, socket *Websocket) bool {
		visitorFunc(socket)
		return false
	})
}
