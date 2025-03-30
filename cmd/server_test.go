package cmd

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestSave_Success(t *testing.T) {
	os.Setenv("DB_USER", "postgre")
	os.Setenv("DB_PASSWORD", "postgre")
	os.Setenv("DB_HOST", "postgre")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "giant_files")

	server := NewStreamServer("localhost", "3000")
	if server == nil {
		t.Fatal("No se pudo inicializar el servidor")
	}

	data := DataEntity{
		FileName: "testfile.txt",
		Total:    42.5,
		Time:     1680000000,
	}

	id, err := server.Save(data)
	if err != nil {
		t.Fatalf("Error al guardar el dato: %v", err)
	}
	if id == 0 {
		t.Errorf("ID insertado no debería ser 0")
	}
}
