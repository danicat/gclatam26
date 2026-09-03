# Gopher Budokai: Panic & Recover!

Jogo de luta de arena 2D desenvolvido para a **GopherCon LATAM 2026 Mini Game Jam**, inspirado na jogabilidade aérea ágil de **Dragon Ball Z: Budokai Tenkaichi 3 (BT3)** e criado estritamente ao redor do tema oficial da jam: **"Panic!!! (& recover?)"**.

---

## Autores

* **Gabriel Donadel** ([@gabrieldonadel](https://github.com/gabrieldonadel))

---

## Tema da Jam: "Panic!!! (& recover?)"

No calor da batalha, receber combos de socos rápidos, ser encurralado por chuvas de *Ki Blasts* ou esgotar sua energia leva o guerreiro ao **PANIC STATE**:
* **PANIC!**: A visão é tomada por um alarme vermelho pulsante, a defesa é reduzida em 50% e o adversário ganha abertura para desferir seu ataque final!
* **(& RECOVER?)**: Em pânico ou arremessado longe, o jogador deve apertar repetidamente `Espaço` (ou `Shift`) para encher a barra de esforço de recuperação e detonar uma colossal **Explosive Kiai Wave**, expulsando o oponente e recuperando a postura instantaneamente!
* **Z-Counter Recover**: Acertar a esquiva com `Shift` no instante exato de um golpe inimigo realiza um teletransporte instantâneo (*Instant Transmission*) atrás do adversário!

---

## Recursos & Destaques Técnicos

* **Trilha Sonora Adaptativa de 6 Canais Gerada por IA (`procedural-composer`)**:
  * Composição polifônica inspirada nos temas clássicos de BT3 ("The Meteor" / "Super Survivor") a 145 BPM em D menor.
  * Transição dinâmica: quando o lutador entra em **PANIC!**, a música se transforma em tempo acelerado com pulso de alarme de emergência e arpejos frenéticos!
  * Zero arquivos de áudio externos: 100% sintetizado em código PCM 16-bit 44.1kHz no boot com zero latência.
* **Mecânicas de Budokai Tenkaichi 3**:
  * Voo livre em 8 direções na arena.
  * **Dragon Dash**: Investida em alta velocidade com aura brilhante.
  * **Ki Charge**: Concentração de energia com partículas; atinge modo **Sparking!** no 100%.
  * **Rush Combos**: Sequência de golpes de artes marciais que termina em arremesso supersônico.
  * **Ki Blasts**: Metralhadora de projéteis de energia.
  * **Super Beam (Kamehameha / Final Flash)** e **Beam Clash** interativo!
* **Gráficos Procedurais Retro-HD**:
  * Silhuetas e trajes clássicos (Goku Super Saiyan vs Vegeta Saiyan Armor).
  * Pool pré-alocado de 600 partículas para zero alocações na renderização.
  * Resolução virtual 640x360 com escala proporcional e suporte a tela cheia.

---

## Controles

| Ação | Tecla | Descrição |
| :--- | :--- | :--- |
| **Mover / Voar** | `W`, `A`, `S`, `D` ou Setas | Voo em 8 direções na arena |
| **Golpe / Combo** | `J` | Sequência de golpes de curta distância |
| **Ki Blast** | `K` | Rajada de esferas de energia |
| **Dragon Dash** | `L` | Voo supersônico com rastro de velocidade |
| **Carregar Ki** | `Espaço` (Segurar) | Concentra Ki (atinge modo *Sparking!* no máximo) |
| **RECOVER!** | `Espaço` ou `Shift` (Mashing) | **Em Pânico: Pressione rápido para recuperar com Kiai!** |
| **Super Beam** | `I` ou `U` | Dispara o raio de energia colossal (Kamehameha) |
| **Teletransporte / Esquiva** | `Shift` | Teletransporta instantaneamente para trás do oponente |
| **Tela Cheia** | `F11` | Alterna modo tela cheia |
| **Pausar / Voltar** | `Esc` | Retorna ao menu na tela de fim de jogo |

---

## Instruções de Execução

### Pré-requisitos
* Go 1.26 ou superior instalado.

### Como Jogar
1. Navegue até a pasta do jogo:
   ```bash
   cd parte3-mini-game-jam/gopher-budokai-panic
   ```
2. Execute diretamente com:
   ```bash
   go run .
   ```
3. Pressione `Espaço` ou `Enter` na tela de título para iniciar a luta!
