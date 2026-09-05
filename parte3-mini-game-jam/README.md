# GopherCon LATAM 2026 Mini Game Jam

Hackathon de desenvolvimento de jogos 2D em Go com Ebitengine v2 e Antigravity.

## Regras da Jam

* Duração: 90 minutos (+30 min de apresentações)
* Formação: Individual ou times
* Tecnologias: Go 1.25+, Ebitengine v2, Antigravity (`agy`), Nano Banana, Lyria 3
* Tema: A ser anunciado no kickoff

## Pré-requisitos do Ambiente

O Ebitengine v2 usa cgo no Linux: depende do ALSA (áudio) e das bibliotecas de
OpenGL/X11 em tempo de **compilação**. Sem elas o `go build` falha antes mesmo do
jogo abrir. Instale antes do kickoff.

**Debian / Ubuntu**
```bash
sudo apt install pkg-config libasound2-dev libgl1-mesa-dev \
  libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libxxf86vm-dev
```

**Fedora**
```bash
sudo dnf install pkgconf-pkg-config alsa-lib-devel mesa-libGL-devel \
  libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel
```

**Arch**
```bash
sudo pacman -S pkgconf alsa-lib mesa libxrandr libxcursor libxinerama libxi
```

**macOS e Windows**: não precisam de nada além do Go — o Ebitengine usa os
drivers nativos do sistema.

Para validar o ambiente antes da jam:
```bash
cd parte3-mini-game-jam/panic-recover && go build ./...
```

## Fluxo de Participação (Fork + PR)

1. Fazer Fork deste repositorio.
2. Clonar localmente:
   ```bash
   git clone https://github.com/SEU-USUARIO/gclatam26.git
   cd gclatam26/parte3-mini-game-jam
   ```
3. Criar pasta do jogo:
   ```bash
   mkdir -p nome-do-jogo
   cd nome-do-jogo
   go mod init nome-do-jogo
   go get github.com/hajimehoshi/ebitengine/v2@latest
   ```
4. Desenvolver com `agy` utilizando as skills em `.agents/skills/`.
5. Criar `README.md` dentro de `parte3-mini-game-jam/nome-do-jogo/` com:
   * Título
   * Autores
   * Controles
   * Instruções de execução
6. Submeter via Pull Request contra a branch `main` com o titulo: `[Game Jam] Nome do Jogo - @autor`.
