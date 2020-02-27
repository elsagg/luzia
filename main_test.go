package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/elsagg/luzia/app"
	"github.com/elsagg/luzia/datacell"
	"github.com/go-testfixtures/testfixtures"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis      *bufconn.Listener
	db       *sql.DB
	fixtures *testfixtures.Context
)

func init() {
	if appEnv != "production" {
		godotenv.Load()
	}
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	datacell.RegisterDataCellServiceServer(s, &app.Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestMain(m *testing.M) {
	var err error

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1)/users?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS users;`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			added_id bigint(20) NOT NULL AUTO_INCREMENT,
			row_key varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
			column_key varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
			body blob,
			ref_key bigint(20) DEFAULT NULL,
			created_at datetime DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (added_id)
		  ) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	if err != nil {
		log.Fatal(err)
	}

	testfixtures.SkipDatabaseNameCheck(true)

	fixtures, err = testfixtures.NewFolder(db, &testfixtures.MySQL{}, "testdata/fixtures")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestGetDataCell(t *testing.T) {
	prepareTestDatabase()

	tests := []struct {
		request *datacell.GetDataCellRequest
		want    *datacell.GetDataCellResponse
	}{
		{
			request: &datacell.GetDataCellRequest{
				DataSource: "users",
				RowKey:     "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey:  "BASIC_INFO",
				RefKey:     2,
			},
			want: &datacell.GetDataCellResponse{
				AddedID:   "3",
				RowKey:    "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey: "BASIC_INFO",
				RefKey:    2,
				Body:      "{\"first_name\":\"Zelda\",\"last_name\":\"Princess\"}",
				CreatedAt: "2018-11-27T16:38:13Z",
			},
		},
		{
			request: &datacell.GetDataCellRequest{
				DataSource: "users",
				RowKey:     "185133a4-71e4-4595-8f09-3ffe416064fc",
				ColumnKey:  "ADDRESS",
				RefKey:     1,
			},
			want: &datacell.GetDataCellResponse{
				AddedID:   "6",
				RowKey:    "185133a4-71e4-4595-8f09-3ffe416064fc",
				ColumnKey: "ADDRESS",
				RefKey:    1,
				Body:      "{\"street\":\"Sunset Boullevard\",\"zip\":105245621}",
				CreatedAt: "2018-11-26T15:14:38Z",
			},
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := datacell.NewDataCellServiceClient(conn)

	for _, tt := range tests {
		resp, err := client.GetDataCell(ctx, tt.request)
		if err != nil {
			t.Fatalf("TestGetDataCell(%v) got unexpected error", err)
		}

		if !reflect.DeepEqual(resp, tt.want) {
			t.Errorf("TestGetDataCell(%v) error: expected %v, got %v", tt.request.GetRowKey(), tt.want, resp)
		}
	}
}

func TestGetLatestDataCell(t *testing.T) {
	prepareTestDatabase()

	tests := []struct {
		request *datacell.GetLatestDataCellRequest
		want    *datacell.GetDataCellResponse
	}{
		{
			request: &datacell.GetLatestDataCellRequest{
				DataSource: "users",
				RowKey:     "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey:  "BASIC_INFO",
			},
			want: &datacell.GetDataCellResponse{
				AddedID:   "3",
				RowKey:    "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey: "BASIC_INFO",
				RefKey:    2,
				Body:      "{\"first_name\":\"Zelda\",\"last_name\":\"Princess\"}",
				CreatedAt: "2018-11-27T16:38:13Z",
			},
		},
		{
			request: &datacell.GetLatestDataCellRequest{
				DataSource: "users",
				RowKey:     "185133a4-71e4-4595-8f09-3ffe416064fc",
				ColumnKey:  "ADDRESS",
			},
			want: &datacell.GetDataCellResponse{
				AddedID:   "7",
				RowKey:    "185133a4-71e4-4595-8f09-3ffe416064fc",
				ColumnKey: "ADDRESS",
				RefKey:    2,
				Body:      "{\"street\":\"Somewhere Inhell\",\"zip\":256255254}",
				CreatedAt: "2018-11-26T15:18:38Z",
			},
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := datacell.NewDataCellServiceClient(conn)

	for _, tt := range tests {
		resp, err := client.GetLatestDataCell(ctx, tt.request)
		if err != nil {
			t.Fatalf("TestGetDataCell(%v) got unexpected error", err)
		}

		if !reflect.DeepEqual(resp, tt.want) {
			t.Errorf("TestGetDataCell(%v) error: expected %v, got %v", tt.request.GetRowKey(), tt.want, resp)
		}
	}
}

func TestPutDataCell(t *testing.T) {
	prepareTestDatabase()

	tests := []struct {
		request *datacell.PutDataCellRequest
		want    *datacell.GetDataCellResponse
	}{
		{
			request: &datacell.PutDataCellRequest{
				DataSource: "users",
				RowKey:     "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey:  "BASIC_INFO",
				Body:       "{\"first_name\":\"Link\",\"last_name\":\"Hero\"}",
			},
			want: &datacell.GetDataCellResponse{
				AddedID:   "9",
				RowKey:    "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey: "BASIC_INFO",
				RefKey:    3,
				Body:      "{\"first_name\":\"Link\",\"last_name\":\"Hero\"}",
			},
		},
		{
			request: &datacell.PutDataCellRequest{
				DataSource: "users",
				RowKey:     "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey:  "ADDRESS",
				Body:       "{\"street\":\"Third Avenue\",\"zip\":256546854}",
			},
			want: &datacell.GetDataCellResponse{
				AddedID:   "10",
				RowKey:    "79440bc9-f928-42ee-970e-17958d2ceb5b",
				ColumnKey: "ADDRESS",
				RefKey:    2,
				Body:      "{\"street\":\"Third Avenue\",\"zip\":256546854}",
			},
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := datacell.NewDataCellServiceClient(conn)

	for _, tt := range tests {
		resp, err := client.PutDataCell(ctx, tt.request)
		if err != nil {
			t.Fatalf("TestGetDataCell(%v) got unexpected error", err)
		}

		if resp.RefKey != tt.want.RefKey {
			t.Errorf("TestGetDataCell(%v) error: expected %v, got %v", tt.request.GetRowKey(), tt.want, resp)
		}
	}
}
