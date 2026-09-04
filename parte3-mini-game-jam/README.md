# GopherCon LATAM 2026 Mini Game Jam

Hackathon de desenvolvimento de jogos 2D em Go com Ebitengine v2 e Antigravity.

## Regras da Jam

* Duração: 90 minutos (+30 min de apresentações)
* Formação: Individual ou times
* Tecnologias: Go 1.26+, Ebitengine v2, Antigravity (`agy`), Nano Banana, Lyria 3
* Tema: A ser anunciado no kickoff

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
