package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type variantConfig struct {
	Width   int
	Bitrate string
}

var variants = map[string]variantConfig{
	"4K":    {Width: 3840, Bitrate: "8000k"},
	"1440p": {Width: 2560, Bitrate: "6000k"},
	"1080p": {Width: 1920, Bitrate: "4000k"},
	"720p":  {Width: 1280, Bitrate: "2500k"},
	"480p":  {Width: 854, Bitrate: "1000k"},
	"360p":  {Width: 640, Bitrate: "500k"},
	"240p":  {Width: 426, Bitrate: "300k"},
}

// ffmpegExec existe separado para facilitar testes/mocks.
var ffmpegExec = func(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "ffmpeg", args...)
}

type FFmpegTranscoder struct{}

func New() *FFmpegTranscoder {
	return &FFmpegTranscoder{}
}

// Transcode gera as variações HLS em 3840p, 1440p, 1080p, 720p, 480p, 360p e 240p.
func (t *FFmpegTranscoder) Transcode(
	ctx context.Context,
	videoID,
	inputPath string,
	resolutions []string,
) error {
	threads := os.Getenv("FFMPEG_THREADS")
	if threads == "" {
		threads = "2" // limitado para evitar processador sobrecarregado
	}

	// Verifica se o arquivo de entrada existe antes de iniciar o ffmpeg.
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Pasta base onde os arquivos HLS vão ser gerados.
	outputDir := filepath.Join("/data/videos", videoID)

	// Valida resoluções e cria diretórios de saída.
	for _, res := range resolutions {
		if _, exists := variants[res]; !exists {
			return fmt.Errorf("unknown resolution %q", res)
		}

		if err := os.MkdirAll(filepath.Join(outputDir, res), 0755); err != nil {
			return fmt.Errorf("failed to create dir %s: %w", res, err)
		}
	}

	args := []string{
		"-hide_banner",
		"-y",
		"-i", inputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-preset", "ultrafast",
		"-threads", threads,
	}

	// Mapeia streams de vídeo e áudio para cada resolução.
	for range resolutions {
		args = append(args,
			"-map", "0:v:0",
			"-map", "0:a:0?",
		)
	}

	// Define resolução e bitrate de cada variante.
	for i, res := range resolutions {
		cfg, ok := variants[res]
		if !ok {
			return fmt.Errorf("unknown resolution %q", res)
		}

		args = append(args,
			fmt.Sprintf("-filter:v:%d", i),
			fmt.Sprintf("scale=%d:-2", cfg.Width),
			fmt.Sprintf("-b:v:%d", i),
			cfg.Bitrate,
		)
	}

	// Monta o var_stream_map.
	var streamMap []string
	for i, res := range resolutions {
		if _, ok := variants[res]; !ok {
			return fmt.Errorf("unknown resolution %q", res)
		}

		streamMap = append(
			streamMap,
			fmt.Sprintf("v:%d,a:%d,name:%s", i, i, res),
		)
	}

	args = append(args,
		"-var_stream_map", strings.Join(streamMap, " "),
		"-master_pl_name", "master.m3u8",
		"-f", "hls",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputDir, "%v", "seg%03d.ts"),
		filepath.Join(outputDir, "%v", "index.m3u8"),
		"-progress", "pipe:2",
		"-nostats",
	)

	log.Printf("executando ffmpeg: ffmpeg %s", strings.Join(args, " "))

	cmd := ffmpegExec(ctx, args...)

	// Captura stderr para ler progresso e erros do ffmpeg.
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	// Inicia o processo antes de começar a ler a saída.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start ffmpeg: %w", err)
	}

	scanErrCh := make(chan error, 1)

	// Lê a saída do ffmpeg em tempo real.
	go func() {
		scanner := bufio.NewScanner(stderr)

		// O ffmpeg pode gerar linhas grandes; aumenta o buffer padrão.
		scanner.Buffer(make([]byte, 1024), 1024*1024)

		for scanner.Scan() {
			log.Printf("[ffmpeg] %s", scanner.Text())
		}

		scanErrCh <- scanner.Err()
	}()

	// Espera o processo terminar.
	waitErr := cmd.Wait()

	// Verifica erro do scanner.
	if scanErr := <-scanErrCh; scanErr != nil {
		return fmt.Errorf("reading ffmpeg output: %w", scanErr)
	}

	if waitErr != nil {
		return fmt.Errorf(
			"ffmpeg failed for video %s: %w",
			videoID,
			waitErr,
		)
	}

	log.Printf("transcoding finalizado: %s", videoID)

	return nil
}
