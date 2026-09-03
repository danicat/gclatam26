# GopherCon LATAM 2026 Mini Game Jam: Panic Invaders

> **Tema Oficial**: *Panic Invaders: In `recover()` We Trust*  
> **Spin-Off Oficial**: *A Batalha pelo Palco Principal da GopherCon LATAM 2026*  
> **Gênero**: Arcade Shoot 'em up 2D (Space Invaders Style)  
> **Engine**: [Ebitengine v2](https://ebitengine.org) (`github.com/hajimehoshi/ebiten/v2`)  
> **Ferramental**: Antigravity (`agy`), Nano Banana 2 Lite / Gemini Image, Lyria 3 / Procedural Audio  

---

## 1. O Tema: Panic vs. Recover

### 1.1 O Spin-Off da GopherCon LATAM 2026
O ano é 2026. Estamos no auditório principal da **GopherCon LATAM**! Centenas de desenvolvedores e entusiastas de Go estão a postos, os telões estão acesos e a contagem regressiva para a aguardada Keynote de Abertura chega a zero.

De repente, os monitores piscam: um commit não revisado em produção disparou uma avalanche incontrolável de goroutines desgovernadas! Uma tempestade catastrófica de `panic()` invadiu os servidores da conferência, ameaçando derrubar a transmissão ao vivo, o sistema de áudio e a infraestrutura do evento.

O processo principal está à beira do colapso fatal (`exit status 2: runtime panicked`).

A única esperança é o lendário **Gopher Defensor**, pilotando o módulo de resgate armado com a última linha de defesa do Go: a função **`recover()`**, operando dentro de uma clausura de emergência protegida por **`defer`**.

Neste jogo:
- **O Herói**: O lendário **Gopher Azul** (saudável, heroico e bom), armado com a sub-rotina **`recover()`**, posicionado na base da Call Stack dos servidores da GopherCon LATAM.
- **Os Inimigos**: As hordas de **Bad Gophers** (vermelhos-escuros, corrompidos e com chifres), trazendo ondas devastadoras de **`panic()`** descendo em bloco para corromper a transmissão da conferência.
- **A Missão**: Interceptar e neutralizar cada exceção antes que alcancem o `main()`, garantindo **99.999% de SLA** e mantendo a GopherCon LATAM 2026 no ar!

---

## 2. Visão Geral do Jogo (Space Invaders Style)

```text
========================================================================
 GOPHERCON LATAM 2026 | SCORE: 004200 | UPTIME: 99.999% | LIVES: ♥ ♥ ♥
========================================================================

       [panic: nil pointer]           [panic: nil pointer]
    [panic: index out of range]    [panic: index out of range]
         [panic: division by zero]    [panic: division by zero]
                     ↓↓↓ DESCE EM BLOCO PELA STACK ↓↓↓

   ------------------------------------------------------------------
      [defer: log.Flush]      [defer: Close()]      [defer: Unlock()]
      (Barreira #1)           (Barreira #2)         (Barreira #3)
   ------------------------------------------------------------------

                                    ▲
                                 [recover()]  <-- GOPHER HERO (A / D + Espaço)
========================================================================
```

### 2.1 O Herói (`recover()`)
- Posicionado na base da tela (`bottom`), com movimentação horizontal (`Left`/`Right` ou `A`/`D`).
- Dispara rajadas verticais de `recover()` (raios azuis ciano característicos do mascote Go) com a tecla `Espaço`.
- Ao atingir um `panic()`, o erro é neutralizado, exibindo logs pontuais (`"recovered: err != nil"`) e somando pontos de estabilidade de infraestrutura.

### 2.2 Os Inimigos: Bad Gophers em Vermelho-Escuro (`panic()`)
Uma matriz invasora de Bad Gophers corrompidos que se move em uníssono de um lado para o outro da Call Stack, descendo um degrau a cada batida de parede:

| Inimigo | Comportamento | Pontos |
| :--- | :--- | :--- |
| **Bad Gopher Scout: `panic("nil pointer")`** | Fileira superior (vermelho-escuro com chifres). Dispara ponteiros corrompidos para baixo. | 30 pts |
| **Bad Gopher Spiky: `panic("index out of range")`** | Fileira intermediária (olhos vermelhos incandescentes). Movimentação veloz e errática. | 20 pts |
| **Bad Gopher Brute: `panic("integer divide by zero")`** | Fileira inferior (mandíbulas afiadas). Frequência elevada de disparos de erro. | 10 pts |
| **Bad Gopher Drone: `panic("wifi connection dropped")`** | Drone dark-red que voa rapidamente pelo topo do auditório em intervalos surpresa. | 100-300 pts |
| **Boss Wave 3: Corrupted Mega-Bad-Gopher (`panic("deadlock")`)** | O temido boss titânico da transmissão ao vivo! Exige dezenas de `recover()` para ser quebrado. | 1000 pts |

### 2.3 Barreiras de Defesa (`defer` blocks)
- Três barreiras protetoras representando rotinas críticas de encerramento (`defer log.Flush()`, `defer db.Close()`, `defer mu.Unlock()`).
- Bloqueiam disparos de ambos os lados, sofrendo dano destrutivo pixel-a-pixel até evaporarem.

### 2.4 Power-ups Temáticos da Conferência & Go
Destruir inimigos especiais ou o UFO do Wi-Fi pode liberar crachás e drops:
- 🛡️ **`sync.Mutex`**: Escudo temporário contra colisões e tiros inimigos.
- ⚡ **`chan struct{}`**: Disparo triplo em leque de canais não-bloqueantes.
- ⏱️ **`context.WithTimeout`**: Congela/desacelera a descida dos bugs por 5 segundos.
- ☕ **`Café do Coffee Break` / `GopherCon LATAM Badge`**: Velocidade de movimento dobrada e pulso que limpa a tela.

### 2.5 Condições de Vitória e Derrota
- **Derrota (Keynote Abortada)**:
  - O jogador perde todas as 3 vidas (Goroutines); OU
  - A horda de `panic()` atinge a linha base do auditório.
  - *Tela de Game Over*: Telão azul de crash simulando `exit status 2: GopherCon LATAM stream terminated`.
- **Vitória (Keynote Sucesso Absoluto)**:
  - Sobreviver às 3 ondas e neutralizar o Deadlock Boss.
  - *Tela de Vitória*: Palco iluminado, aplausos e a mensagem: `GopherCon LATAM 2026 salva com sucesso! 100% Uptime!`.

---

## 3. Controles

| Ação | Teclado | Gamepad | Mobile / Touch (Opcional) |
| :--- | :--- | :--- | :--- |
| Mover Esquerda | `A` ou `Seta Esquerda` | D-Pad Esquerdo / Analógico | Toque no terço esquerdo |
| Mover Direita | `D` ou `Seta Direita` | D-Pad Direito / Analógico | Toque no terço direito |
| Disparar `recover()` | `Espaço` ou `J` | Botão A / X (South) | Toque na área central |
| Pausar | `Esc` ou `P` | Start | Ícone de Pause |

---

## 4. Stack Tecnológica & Assets

1. **Engine**:
   - [Ebitengine v2](https://ebitengine.org) (`github.com/hajimehoshi/ebiten/v2`).
2. **Visual & Arte**:
   - **Opção A (Procedural)**: Desenho procedural via primitivas Ebitengine (`ebitenutil.DrawRect`, vetores, Kage shaders) com o skill `procedural-art`.
   - **Opção B (GenAI)**: Sprites do Gopher e bugs gerados com Nano Banana 2 Lite / Gemini (`google.golang.org/genai`) utilizando o script da `parte1-gemini`.
3. **Áudio**:
   - **SFX**: Sons retro/chiptune gerados por código procedural (`procedural-composer`).
   - **BGM**: Trilha sonora synthwave sintetizada via Lyria 3 utilizando o microserviço da `parte2-antigravity`.
4. **Agentes & Skills**:
   - Disponíveis em `.agents/skills/`: `ebitengineer`, `game-design`, `procedural-art`, `procedural-composer`, `godoctor`.

---

## 5. Regras da Jam

* **Duração**: 90 minutos (+30 min de apresentações e playtest).
* **Formato**: Individual ou duplas/times.
* **Tema Obrigatório**: `panic` vs `recover` no estilo Space Invaders (com spin-off GopherCon LATAM).
* **Critérios de Votação**:
  1. **Jogabilidade (Game Feel)**: Quão satisfatório e fluido é o combate.
  2. **Aderência ao Tema & Humor Go**: Criatividade no uso das piadas de infraestrutura, conferência e runtime Go.
  3. **Originalidade Audiovisual**: Uso de efeitos sonoros, trilha musical e visual arcade.
  4. **Estabilidade de Código**: Código Go idiomático e sem panics não tratados!

---

## 6. Fluxo de Submissão (Fork + PR)

1. Fazer Fork deste repositório para a sua conta GitHub.
2. Clonar localmente:
   ```bash
   git clone https://github.com/SEU-USUARIO/gclatam26.git
   cd gclatam26/parte3-mini-game-jam
   ```
3. Criar a pasta do seu jogo:
   ```bash
   mkdir -p panic-invaders
   cd panic-invaders
   go mod init panic-invaders
   go get github.com/hajimehoshi/ebiten/v2@latest
   ```
4. Desenvolver o jogo com auxílio do `agy` e das skills locais.
5. Criar um `README.md` dentro de `parte3-mini-game-jam/panic-invaders/` contendo:
   - Nome do Jogo & Autores
   - Instruções de execução (`go run main.go`)
   - Controles e mecânicas implementadas
   - Screenshot ou GIF do jogo
6. Submeter via Pull Request contra a branch `main` do repositório original com o título:  
   `[Game Jam] Nome do Jogo - @autor`
