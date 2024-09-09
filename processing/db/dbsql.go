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
	GetPedidos(ctx context.Context, usuarioId string) ([]domain.PedidoDetalhado, float64, error) 
	DeletarPedidos(ctx context.Context) error
}

func (s *SQLStore) GetItem(ctx context.Context, nome string) (*domain.Item, error) {
	var item domain.Item

	query := `
		SELECT nome, tempo_corte, tempo_grelha, tempo_montagem, tempo_bebida, valor
		FROM item
		WHERE nome = $1`

	row := s.db.QueryRowContext(ctx, query, nome)

	var tempoCorteStr, tempoGrelhaStr, tempoMontagemStr, tempoBebidaStr string
	err := row.Scan(&item.Nome, &tempoCorteStr, &tempoGrelhaStr, &tempoMontagemStr, &tempoBebidaStr, &item.Valor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.New(exceptions.ErrOrderNotFound, err)
		}
		log.Error("Erro ao buscar item no banco de dados: ", err)
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
		return 0, fmt.Errorf("formato de tempo inválido, esperava hh:mm:ss: %s", hms)
	}

	hours := parts[0]
	minutes := parts[1]
	seconds := parts[2]

	durationStr := fmt.Sprintf("%sh%sm%ss", hours, minutes, seconds)
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("erro ao parsear a duração: %w", err)
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
			log.Printf("erro ao escanear o pedido: %v", err)
			continue
		}

		item, err := s.GetItem(ctx, itemId)
		if err != nil {
			log.Printf("erro ao buscar item no banco: %v", err)
			continue
		}

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
		log.Error("erro ao criar pedido: ", err)
		return exceptions.New(exceptions.ErrInternalServer, err)
	}
	return nil
}

func (s *SQLStore) DeletarPedidos(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		TRUNCATE TABLE pedido
		`,
	)
	if err != nil {
		log.Error("erro ao limpae tabela pedido: ", err)
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
		log.Error("erro ao criar pedido: ", err)
		return exceptions.New(exceptions.ErrInternalServer, err)
	}
	return nil
}

func (s *SQLStore) GetPedidos(ctx context.Context, usuarioId string) ([]domain.PedidoDetalhado, float64, error) {
    idsQuery := `
    SELECT item_id
    FROM pedido
    WHERE user_id = $1 AND status != 'Cancelado'
    `

    rows, err := s.db.QueryContext(ctx, idsQuery, usuarioId)
    if err != nil {
        return nil, 0, fmt.Errorf("erro ao obter os IDs dos pedidos: %w", err)
    }
    defer rows.Close()

    var ids []string
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, 0, fmt.Errorf("erro ao escanear o ID do pedido: %w", err)
        }
        ids = append(ids, id)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("erro durante a iteração dos IDs: %w", err)
    }

    if len(ids) == 0 {
        return nil, 0, nil
    }

    detalhesQuery := `
    SELECT nome, valor
    FROM item
    WHERE nome = ANY($1)
    `

    idsArray := "{" + strings.Join(ids, ",") + "}"

    rows, err = s.db.QueryContext(ctx, detalhesQuery, idsArray)
    if err != nil {
        return nil, 0, fmt.Errorf("erro ao obter os detalhes dos pedidos: %w", err)
    }
    defer rows.Close()

    var pedidosDetalhados []domain.PedidoDetalhado
    var valorTotal float64

    for rows.Next() {
        var detalhe domain.PedidoDetalhado
        if err := rows.Scan(&detalhe.Nome, &detalhe.Valor); err != nil {
            return nil, 0, fmt.Errorf("erro ao escanear detalhes do pedido: %w", err)
        }
        pedidosDetalhados = append(pedidosDetalhados, detalhe)
        valorTotal += detalhe.Valor
    }

    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("erro durante a iteração das linhas: %w", err)
    }

    return pedidosDetalhados, valorTotal, nil
}