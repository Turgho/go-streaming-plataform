package usecase

import "errors"

var (
	ErrForbidden = errors.New("acesso negado ao vídeo")
	ErrNotFound  = errors.New("vídeo não encontrado")
)
