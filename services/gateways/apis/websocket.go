package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"

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

	// Make sure that we automatically clean up the registry when connections close.
	socket := Websocket{Conn: conn, ID: connectionID, Options: opts.applyDefaults(), newMessageContext: newMessageContext}
	customOnClose := socket.Options.OnClose
	socket.Options.OnClose = func() {
		sockets.remove(connectionID)
		customOnClose()
	}

	// Make sure that we only have one connection for a given connection id. If the same user can be connected in
	// different places, don't use the connection ID "user.123.web" because if they're logged in w/ 2 different
	// browsers they'll clobber each other. Instead, do "user.123.web.TIMESTAMP" or something like that. They'll
	// each maintain connections, and you can locate them both by doing prefix lookups on "user.123.web"
	if old, ok := sockets.add(connectionID, &socket); ok {
		quiet.Close(old)
	}

	socket.startListening()
	return &socket, nil
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
	// OnClose provides a custom handler that fires when this websocket is closed for any reason.
	OnClose func()
}

func (opts WebsocketOptions) applyDefaults() WebsocketOptions {
	if opts.OnReadText == nil {
		opts.OnReadText = func(ctx context.Context, socket *Websocket, data []byte) {}
	}
	if opts.OnReadBinary == nil {
		opts.OnReadBinary = func(ctx context.Context, socket *Websocket, data []byte) {}
	}
	if opts.OnReadContinuation == nil {
		opts.OnReadContinuation = func(ctx context.Context, socket *Websocket, data []byte) {}
	}
	if opts.OnClose == nil {
		opts.OnClose = func() {}
	}
	return opts
}

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
}

// Active returns true if the underlying connection has NOT been closed yet.
func (socket *Websocket) Active() bool {
	return socket.Conn != nil
}

// WriteText writes a frame of binary data to the client on the other end of the socket.
func (socket *Websocket) Write(data []byte) error {
	if !socket.Active() {
		return fail.Unavailable("socket closed")
	}

	if err := wsutil.WriteServerBinary(socket.Conn, data); err != nil {
		quiet.Close(socket)
		return fmt.Errorf("error writing to websocket: %s: %w", socket.ID, err)
	}
	return nil
}

// WriteClose pushes a "close" frame to the client, letting them know that we want to close up shop.
func (socket *Websocket) WriteClose(data []byte) error {
	if !socket.Active() {
		return fail.Unavailable("socket closed")
	}

	if err := wsutil.WriteServerMessage(socket.Conn, ws.OpClose, data); err != nil {
		quiet.Close(socket)
		return fmt.Errorf("error writing to websocket: %s: %w", socket.ID, err)
	}
	return nil
}

// WriteText writes a frame of text data to the client on the other end of the socket.
func (socket *Websocket) WriteText(data string) error {
	if !socket.Active() {
		return fail.Unavailable("socket closed")
	}

	if err := wsutil.WriteServerText(socket.Conn, []byte(data)); err != nil {
		quiet.Close(socket)
		return fmt.Errorf("error writing to websocket: %s: %w", socket.ID, err)
	}
	return nil
}

// WriteJSON marshals the given object into a JSON string and then writes a text frame to the client.
func (socket *Websocket) WriteJSON(value any) error {
	if !socket.Active() {
		return fail.Unavailable("socket closed")
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

// Close kills the current connection. This will also trigger your OnClose handler.
func (socket *Websocket) Close() error {
	if !socket.Active() {
		return nil
	}

	if socket.Options.Logger != nil {
		socket.Options.Logger.Info("Websocket CLOSED", "websocket_id", socket.ID)
	}

	quiet.Close(socket.Conn)
	socket.Conn = nil // make sure this socket reference is not seen as 'Active' anymore.
	socket.Options.OnClose()
	return nil
}

// startListening fires off a separate goroutine that infinitely loops, attempting to read messages from the client.
// Any incoming messages will be routed to the socket's OnReadText/Binary/etc. handlers. This will exit automatically
// when the socket/connection is closed.
func (socket *Websocket) startListening() {
	if socket.Options.Logger != nil {
		socket.Options.Logger.Info("Websocket OPENED", "websocket_id", socket.ID)
	}

	go func() {
		defer quiet.Close(socket)
		logger := socket.Options.Logger

		for socket.Active() {
			data, op, err := wsutil.ReadClientData(socket.Conn)
			if err != nil {
				logger.Debug("error reading client data, closing connection",
					"error", err,
					"websocket_id", socket.ID,
				)
				break
			}

			switch op {
			case ws.OpClose:
				break
			case ws.OpText:
				socket.Options.OnReadText(socket.newMessageContext(), socket, data)
			case ws.OpBinary:
				socket.Options.OnReadBinary(socket.newMessageContext(), socket, data)
			case ws.OpContinuation:
				socket.Options.OnReadContinuation(socket.newMessageContext(), socket, data)
			case ws.OpPing:
				_ = wsutil.WriteServerMessage(socket.Conn, ws.OpPong, data)
			case ws.OpPong:
				// Ignore... clients shouldn't be sending pongs, anyway.
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
