package api

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Exercise the HTTP route and database/sql query without a production database.
type healthConnector struct {
	fail    bool
	queries []string
}

func (c *healthConnector) Connect(context.Context) (driver.Conn, error) { return healthConn{c}, nil }
func (c *healthConnector) Driver() driver.Driver                        { return healthDriver{} }

type healthDriver struct{}

func (healthDriver) Open(string) (driver.Conn, error) { return nil, errors.New("connector required") }

type healthConn struct{ c *healthConnector }

func (healthConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unexpected prepare") }
func (healthConn) Close() error                        { return nil }
func (healthConn) Begin() (driver.Tx, error) {
	return nil, errors.New("health must not start a transaction")
}
func (c healthConn) QueryContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	c.c.queries = append(c.c.queries, query)
	if _, ok := ctx.Deadline(); !ok {
		return nil, errors.New("missing query timeout")
	}
	if c.c.fail {
		return nil, errors.New("private connection details")
	}
	return &healthRows{}, nil
}

type healthRows struct{ done bool }

func (*healthRows) Columns() []string { return []string{"one"} }
func (*healthRows) Close() error      { return nil }
func (r *healthRows) Next(values []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	values[0] = int64(1)
	return nil
}

func TestHealthRoute(t *testing.T) {
	for _, name := range []string{"healthy", "query failure", "missing database"} {
		t.Run(name, func(t *testing.T) {
			h := &HealthResource{}
			connector := &healthConnector{fail: name == "query failure"}
			if name != "missing database" {
				db := sql.OpenDB(connector)
				t.Cleanup(func() { _ = db.Close() })
				var err error
				h.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{DisableAutomaticPing: true})
				if err != nil {
					t.Fatal(err)
				}
			}
			ws := new(restful.WebService).Path("/").Produces(restful.MIME_JSON)
			h.Register(ws)
			container := restful.NewContainer()
			container.Add(ws)
			response := httptest.NewRecorder()
			container.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/health", nil))
			wantCode, wantStatus := http.StatusServiceUnavailable, "DOWN"
			if name == "healthy" {
				wantCode, wantStatus = http.StatusOK, "UP"
			}
			if response.Code != wantCode {
				t.Fatalf("status %d: %s", response.Code, response.Body.String())
			}
			if response.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("health response must not be cached")
			}
			if !strings.Contains(response.Header().Get("Content-Type"), "application/json") {
				t.Fatal("expected JSON")
			}
			var body HealthResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			component, exists := body.Components["postgres"]
			if body.Status != wantStatus || !exists || component.Status != wantStatus || component.LatencyMS < 0 {
				t.Fatalf("unexpected health: %+v", body)
			}
			if strings.Contains(response.Body.String(), "private") {
				t.Fatal("dependency error leaked")
			}
			if name != "missing database" && (len(connector.queries) != 1 || connector.queries[0] != "SELECT 1") {
				t.Fatalf("unexpected queries: %v", connector.queries)
			}
		})
	}
}
