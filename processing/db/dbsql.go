package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hsxflowers/restaurante-digital/exceptions"
	"github.com/hsxflowers/restaurante-digital/processing/domain"
	"github.com/labstack/gommon/log"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

type RestauranteDatabase interface {
	CreatePedido(ctx context.Context, pedido *domain.Pedido) error
	GetItem(ctx context.Context, itemId string) (*domain.Item, error)
	GetPedidosAnteriores(ctx context.Context, pedidoId string) ([]domain.Pedido, error)
	UpdatePedidoStatus(ctx context.Context, pedidoId string, status string) error
}

func (s *SQLStore) GetItem(ctx context.Context, itemId string) (*domain.Item, error) {
	var item domain.Item

	query := `
		SELECT item_id, nome, tempo_corte, tempo_grelha, tempo_montagem, tempo_bebida, valor
		FROM item
		WHERE item_id = $1`

	row := s.db.QueryRowContext(ctx, query, itemId)

	var tempoCorteStr, tempoGrelhaStr, tempoMontagemStr, tempoBebidaStr string
	err := row.Scan(&item.ItemId, &item.Nome, &tempoCorteStr, &tempoGrelhaStr, &tempoMontagemStr, &tempoBebidaStr, &item.Valor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.New(exceptions.ErrOrderNotFound, err)
		}
		log.Error("Error fetching item from database: ", err)
		return nil, exceptions.New(exceptions.ErrInternalServer, err)
	}

	item.TempoCorte, err = parseInterval(tempoCorteStr)
	if err != nil {
		log.Error("Erro ao converter tempo_corte: ", err)
		return nil, err
	}
	item.TempoGrelha, err = parseInterval(tempoGrelhaStr)
	if err != nil {
		log.Error("Erro ao converter tempo_grelha: ", err)
		return nil, err
	}
	item.TempoMontagem, err = parseInterval(tempoMontagemStr)
	if err != nil {
		log.Error("Erro ao converter tempo_montagem: ", err)
		return nil, err
	}
	item.TempoBebida, err = parseInterval(tempoBebidaStr)
	if err != nil {
		log.Error("Erro ao converter tempo_bebida: ", err)
		return nil, err
	}

	return &item, nil
}

func parseInterval(hms string) (time.Duration, error) {
	parts := strings.Split(hms, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format, expected hh:mm:ss: %s", hms)
	}

	hours := parts[0]
	minutes := parts[1]
	seconds := parts[2]

	durationStr := fmt.Sprintf("%sh%sm%ss", hours, minutes, seconds)
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing duration: %w", err)
	}

	return duration, nil
}

func (s *SQLStore) GetPedidosAnteriores(ctx context.Context, pedidoId string) ([]domain.Pedido, error) {
	query := `
		SELECT item_id, valor, status
		FROM pedido;
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos anteriores: %w", err)
	}
	defer rows.Close()

	var pedidosAnteriores []domain.Pedido

	for rows.Next() {
		var pedido domain.Pedido
		var itemId string

		if err := rows.Scan(
			&itemId,
			&pedido.Valor,
			&pedido.Status,
		); err != nil {
			log.Printf("Erro ao escanear o pedido: %v", err)
			continue
		}

		item, err := s.GetItem(ctx, itemId)
		if err != nil {
			log.Printf("Erro ao buscar item no banco: %v", err)
			continue
		}

		pedido.ItemId = itemId
		pedido.Nome = item.Nome
		pedido.TempoCorte = item.TempoCorte
		pedido.TempoGrelha = item.TempoGrelha
		pedido.TempoMontagem = item.TempoMontagem
		pedido.TempoBebida = item.TempoBebida

		pedidosAnteriores = append(pedidosAnteriores, pedido)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return pedidosAnteriores, nil
}

func (s *SQLStore) CreatePedido(ctx context.Context, pedido *domain.Pedido) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO pedido (pedido_id, user_id, item_id, valor, status) 
		VALUES ($1, $2, $3, $4, $5)`,
		pedido.PedidoId, pedido.UsuarioId, pedido.ItemId, pedido.Valor, pedido.Status,
	)
	if err != nil {
		log.Error("Error creating pedido in database: ", err)
		return exceptions.New(exceptions.ErrInternalServer, err)
	}
	return nil
}

func (s *SQLStore) UpdatePedidoStatus(ctx context.Context, pedidoId string, status string) error {
	query := `
		UPDATE pedido
		SET status = $1
		WHERE pedido_id = $2
	`

	_, err := s.db.ExecContext(ctx, query, status, pedidoId)
	if err != nil {
		log.Error("Error creating pedido in database: ", err)
		return exceptions.New(exceptions.ErrInternalServer, err)
	}
	return nil
}