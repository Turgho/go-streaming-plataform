#!/usr/bin/env python3
"""
Auto-updater do README.md para go-streaming-platform.

Seções atualizadas automaticamente:
  - SERVICES  → lê docker-compose.yml
  - MAKEFILE  → lê Makefile
  - TREE      → gera árvore real do repositório
  - Badge Go  → lê go.mod de cada serviço (pega a mais recente)
"""

import os
import re

README_PATH = "README.md"

# Pastas/arquivos ignorados na árvore
TREE_IGNORE = {
    ".git", ".github", "node_modules", "__pycache__",
    ".idea", "vendor", "pb", ".env",
    # ignora arquivos de vídeo soltos na raiz
}
TREE_IGNORE_EXTENSIONS = {".mp4", ".mkv", ".avi", ".mov", ".sum"}

# Descrições dos serviços (mantidas aqui para não depender de comentários no compose)
SERVICE_DESCRIPTIONS = {
    "user-service":      "Cadastro, login e validação de token JWT",
    "upload-service":    "Upload de vídeos em chunks via client streaming gRPC",
    "transcode-service": "Transcodificação assíncrona de vídeos via NATS + FFmpeg",
}

# Descrições dos comandos do Makefile
MAKE_DESCRIPTIONS = {
    "proto-user":              "Gera código Go do proto do user-service",
    "proto-upload":            "Gera código Go do proto do upload-service",
    "proto-userpb-upload":     "Copia user proto para o upload-service",
    "proto-uploadpb-transcode":"Copia upload proto para o transcode-service",
    "proto-all":               "Gera todos os protos",
    "docker-up":               "Sobe todos os serviços e bancos",
    "docker-down":             "Para os containers",
}

# ── Helpers ───────────────────────────────────────────────────────────────────

def replace_block(content: str, tag: str, new_block: str) -> str:
    """Substitui o conteúdo entre <!-- TAG_START --> e <!-- TAG_END -->."""
    pattern = rf"(<!-- {tag}_START -->).*?(<!-- {tag}_END -->)"
    replacement = rf"\1\n{new_block}\n\2"
    result = re.sub(pattern, replacement, content, flags=re.DOTALL)
    if result == content:
        print(f"[WARN] Marcador {tag}_START/END não encontrado no README.")
    return result

# ── Serviços ──────────────────────────────────────────────────────────────────

def parse_services() -> str:
    try:
        import yaml
    except ImportError:
        print("[ERRO] Instale pyyaml: pip install pyyaml")
        return "_Erro ao gerar tabela de serviços._"

    try:
        with open("docker-compose.yml", encoding="utf-8") as f:
            compose = yaml.safe_load(f)
    except FileNotFoundError:
        print("[WARN] docker-compose.yml não encontrado.")
        return "_docker-compose.yml não encontrado._"

    rows = ["| Serviço | Porta | Descrição |", "|---|---|---|"]

    for name, cfg in (compose.get("services") or {}).items():
        # Ignora serviços de infraestrutura (sem porta exposta ou sem descrição customizada)
        if name not in SERVICE_DESCRIPTIONS:
            continue
        cfg = cfg or {}
        ports = cfg.get("ports", [])
        if ports:
            # pega apenas a porta do container (lado direito de "HOST:CONTAINER")
            port = f"`:{str(ports[0]).split(':')[-1]}`"
        else:
            port = "—"
        desc = SERVICE_DESCRIPTIONS[name]
        rows.append(f"| **{name}** | {port} | {desc} |")

    return "\n".join(rows)

# ── Makefile ──────────────────────────────────────────────────────────────────

def parse_makefile() -> str:
    try:
        with open("Makefile", encoding="utf-8") as f:
            lines = f.readlines()
    except FileNotFoundError:
        print("[WARN] Makefile não encontrado.")
        return "_Makefile não encontrado._"

    targets = []
    for line in lines:
        # Captura targets reais (não variáveis, não includes, não .PHONY)
        match = re.match(r"^([a-zA-Z0-9_-]+)\s*:", line)
        if match:
            name = match.group(1)
            if not name.startswith("."):
                targets.append(name)

    if not targets:
        return "_Nenhum comando encontrado no Makefile._"

    lines_out = ["```bash"]
    for cmd in targets:
        desc = MAKE_DESCRIPTIONS.get(cmd, "")
        if desc:
            lines_out.append(f"make {cmd:<30} # {desc}")
        else:
            lines_out.append(f"make {cmd}")
    lines_out.append("```")
    return "\n".join(lines_out)

# ── Árvore de diretórios ──────────────────────────────────────────────────────

def generate_tree() -> str:
    def should_ignore(name: str) -> bool:
        if name in TREE_IGNORE:
            return True
        _, ext = os.path.splitext(name)
        return ext in TREE_IGNORE_EXTENSIONS

    def walk(path: str, prefix: str = "", depth: int = 0) -> list:
        if depth > 5:
            return []
        try:
            entries = sorted(os.scandir(path), key=lambda e: (e.is_file(), e.name))
        except PermissionError:
            return []

        entries = [e for e in entries if not should_ignore(e.name)]
        lines = []
        for i, entry in enumerate(entries):
            is_last = i == len(entries) - 1
            connector = "└── " if is_last else "├── "
            lines.append(f"{prefix}{connector}{entry.name}")
            if entry.is_dir():
                extension = "    " if is_last else "│   "
                lines.extend(walk(entry.path, prefix + extension, depth + 1))
        return lines

    root_name = os.path.basename(os.path.abspath("."))
    tree_lines = ["```", f"{root_name}/"] + walk(".") + ["```"]
    return "\n".join(tree_lines)

# ── Versão do Go ──────────────────────────────────────────────────────────────

def find_go_version() -> str | None:
    """Percorre todos os go.mod do monorepo e retorna a versão mais alta."""
    import subprocess

    result = subprocess.run(
        ["find", ".", "-name", "go.mod",
         "-not", "-path", "*/.git/*",
         "-not", "-path", "*/vendor/*"],
        capture_output=True, text=True
    )
    versions = []
    for mod_path in result.stdout.strip().splitlines():
        try:
            with open(mod_path, encoding="utf-8") as f:
                for line in f:
                    m = re.match(r"^go\s+(\d+\.\d+(?:\.\d+)?)", line)
                    if m:
                        versions.append(m.group(1))
                        break
        except Exception:
            continue

    if not versions:
        return None

    def version_key(v):
        return tuple(int(x) for x in v.split("."))

    return max(versions, key=version_key)

def update_go_badge(content: str) -> str:
    version = find_go_version()
    if not version:
        print("[WARN] Versão do Go não encontrada.")
        return content

    updated = re.sub(
        r"(img\.shields\.io/badge/Go-)[\d.]+(\+)",
        rf"\g<1>{version}\2",
        content
    )
    if updated != content:
        print(f"[OK] Badge do Go atualizado para {version}")
    return updated

# ── Stack: tecnologias novas (NATS) ──────────────────────────────────────────

def update_stack_table(content: str) -> str:
    """
    Garante que NATS aparece na tabela de Stack se o compose usar nats.
    Insere a linha apenas se ainda não estiver presente.
    """
    try:
        with open("docker-compose.yml", encoding="utf-8") as f:
            raw = f.read()
    except FileNotFoundError:
        return content

    if "nats" in raw and "NATS" not in content:
        nats_row = "| **NATS**            | Mensageria assíncrona entre serviços |"
        # Insere antes do fechamento da tabela de stack (antes do primeiro ---)
        # Encontra a tabela de Stack e adiciona a linha
        pattern = r"(\| \*\*GitHub Actions\*\*.*?\|.*?\|)"
        replacement = r"\1\n" + nats_row
        updated = re.sub(pattern, replacement, content)
        if updated != content:
            print("[OK] NATS adicionado à tabela de Stack.")
        return updated

    return content

# ── Main ──────────────────────────────────────────────────────────────────────

def main():
    with open(README_PATH, encoding="utf-8") as f:
        content = f.read()

    print("🔄 Atualizando README.md...")

    content = replace_block(content, "SERVICES", parse_services())
    print("[OK] Seção SERVICES atualizada.")

    content = replace_block(content, "MAKEFILE", parse_makefile())
    print("[OK] Seção MAKEFILE atualizada.")

    content = replace_block(content, "TREE", generate_tree())
    print("[OK] Seção TREE atualizada.")

    content = update_go_badge(content)
    content = update_stack_table(content)

    with open(README_PATH, "w", encoding="utf-8") as f:
        f.write(content)

    print("✅ README.md atualizado com sucesso!")

if __name__ == "__main__":
    main()