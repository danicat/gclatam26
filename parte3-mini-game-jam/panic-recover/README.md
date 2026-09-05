# Panic & Recover: Eldritch Gopher

Jogo desenvolvido para a **GopherCon LATAM 2026 Mini Game Jam**.

* **Tema**: *Panic & Recover*
* **Gênero**: Grid-Based Turn-Budget Puzzle (estilo *Helltaker* / *Sokoban*)
* **Autores**: Julio Tsutsui (@JulioTsutsui) & Antigravity

---

## 🎮 Sobre o Jogo

O intrépido mascote Gopher explora catacumbas cósmicas para **recuperar artefatos ancestrais** antes que a loucura o consuma!

* **Panic (Pânico)**: Cada movimento gasta um turno do relógio mental. Ao ultrapassar **80% de pânico**, a realidade se distorce com aberração cromática visual e batimentos cardíacos acelerados. Aos 100%, o Gopher desmaia de exaustão mental!
* **Recover (Recuperar)**: Colete os **Relógios / Cafés de Sanidade** espalhados pelas salas para reverter o relógio, diminuir o pânico abaixo de 80% e recuperar a clareza mental antes de alcançar o **Artefato Eldritch**.
* **Buracos e Pedras**: Empurre monólitos de pedra ancestrais para dentro de fossas cósmicas para criar pontes seguras e abrir caminho.

---

## 🕹️ Controles

| Ação | Teclas |
| :--- | :--- |
| **Mover / Empurrar** | `W`, `A`, `S`, `D` ou `Setas do Teclado` |
| **Reiniciar Sala (Instantâneo)** | `R` |
| **Avançar / Confirmar** | `Espaço` ou `Enter` |
| **Tela Cheia** | `F11` |

---

## 🚀 Instruções de Execução

### Pré-requisitos
* Go 1.25+
* No Linux, as dependências de sistema do Ebitengine v2 (`pkg-config`, ALSA, OpenGL/X11).
  Comandos de instalação por distro em [../README.md](../README.md#pré-requisitos-do-ambiente).
  macOS e Windows não precisam de nada além do Go.

### Executando
```bash
cd panic-recover
go run .
```

### Executando os Testes
```bash
go test -v ./internal/game/...
```
