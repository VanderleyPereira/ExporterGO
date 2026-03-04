package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Config armazena as configurações de conexão
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// DB é a instância global de conexão
var DB *sql.DB

// Connect abre a conexão com o MySQL
func Connect(cfg Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	DB = db
	fmt.Printf("✅ Conectado ao banco: %s\n", cfg.DBName)
	return nil
}

// ListTables retorna todas as tabelas do banco conectado
func ListTables() ([]string, error) {
	rows, err := DB.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("erro ao listar tabelas: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, nil
}

// QueryTable busca todos os dados de uma tabela dinamicamente
// Retorna cabeçalhos e linhas sem precisar de struct fixa
func QueryTable(table string) ([]string, [][]string, error) {
	rows, err := DB.Query(fmt.Sprintf("SELECT * FROM `%s`", table))
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao consultar tabela %s: %w", table, err)
	}
	defer rows.Close()

	// Captura os nomes das colunas
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var data [][]string

	for rows.Next() {
		// Cria um slice de interface{} para receber os valores
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, nil, err
		}

		// Converte tudo para string para exportação universal
		row := make([]string, len(cols))
		for i, val := range values {
			switch v := val.(type) {
			case nil:
				row[i] = "NULL"
			case []byte:
				row[i] = string(v)
			default:
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		data = append(data, row)
	}

	return cols, data, nil
}

// QueryCustom executa uma query customizada do usuário
func QueryCustom(query string) ([]string, [][]string, error) {
	rows, err := DB.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("erro na query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var data [][]string
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, nil, err
		}
		row := make([]string, len(cols))
		for i, val := range values {
			switch v := val.(type) {
			case nil:
				row[i] = "NULL"
			case []byte:
				row[i] = string(v)
			default:
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		data = append(data, row)
	}

	return cols, data, nil
}
