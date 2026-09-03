# Exploding Kitchens (Cozinhas Explosivas)

> **Evento**: GopherCon LATAM 2026 Mini Game Jam  
> **Tema**: Panic and Recover (Pânico e Recuperação)  
> **Inspiração**: *Exploding Kittens* × *Overcooked* / Catástrofe Culinária  
> **Tecnologias**: Go 1.26+, Ebitengine v2 (`github.com/hajimehoshi/ebiten/v2`)  
> **Autores**: Joubert RedRat & Antigravity AI Pair Programmer  

---

## 1. Sobre o Jogo

Em **Exploding Kitchens**, você assume o papel de um chef em apuros durante o horário de pico mais caótico do mundo! Eletrodomésticos experimentais de alta voltagem entram em superaquecimento contínuo, e para piorar, um gato travesso passeia livremente pelas bancadas, girando botões para o fogo máximo e ameaçando explodir o restaurante inteiro.

O seu objetivo é simples: **sobreviver a um turno frenético de 2 minutos** controlando a **Barra de Caos**, desarmando aparelhos à beira da detonação e consertando os estragos antes que o restaurante vá pelos ares!

---

## 2. Como Jogar (Guia Passo a Passo)

### 2.1 O Ciclo Principal de Jogo
1. **Fique de Olho nos Eletrodomésticos**:
   * As estações de cozinha (Panela de Pressão, Fogão, Fritadeira e Micro-ondas) possuem uma barra de perigo flutuante.
   * **Verde (Cozinhando)**: Operação normal.
   * **Amarelo (Alerta)**: Começa a soltar vapor e emitir bipes de aviso.
   * **Vermelho Piscante (PÂNICO!)**: Sirenes tocam, chamas e fumaça sobem. A explosão é iminente!
2. **Pegue a Ferramenta Certa**:
   * Corra até a prateleira de ferramentas na parte inferior da cozinha e pressione `Espaço` ou `E` para equipar a ferramenta necessária.
3. **Desarme a Estação**:
   * Aproxime-se do aparelho em perigo com a ferramenta correta e aperte `Espaço` ou `E` para desarmá-lo e zerar o perigo.
4. **Espante o Gato Travesso**:
   * O gato laranja pula nas bancadas dos eletrodomésticos e **acelera o tempo de explosão em 2.5x**! Aproxime-se dele e aperte `Espaço` (ou simplesmente encoste nele) para espantá-lo de volta para o tapetinho rosa.
5. **Conserte com a Chave Inglesa**:
   * Se um aparelho explodir, ele ficará destruído e soltando fumaça preta. Pegue a **Chave Inglesa** e vá até ele para consertá-lo.

---

## 3. Matriz de Eletrodomésticos e Ferramentas

| Eletrodoméstico | Comportamento de Pânico | Ferramenta Necessária | Efeito da Recuperação |
| :--- | :--- | :--- | :--- |
| **Panela de Pressão** *(Superior Esquerda)* | Acumula pressão em alta velocidade e apita. | **Balde de Gelo** *(ou mãos vazias)* | Resfria a temperatura e alivia o vapor com segurança. |
| **Fogão com Queimadores** *(Superior Centro-Esq)* | Chamas sobem e panelas transbordam óleo. | **Extintor** ou **Balde de Gelo** | Apaga as chamas e resfria os queimadores. |
| **Fritadeira de Óleo** *(Superior Centro-Dir)* | Óleo ferve violentamente e entra em combustão. | **Extintor de Incêndio** | Sufoca o fogo de óleo e reinicia a fritura. |
| **Micro-ondas Bomba** *(Superior Direita)* | Timer digital regressivo em contagem acelerada. | **Balde de Gelo** *(ou mãos vazias)* | Resfria os circuitos e reinicia o temporizador. |
| **Estação Explodida** *(Destruída/Carbonizada)* | Fica inutilizada e soltando faíscas após detonar. | **Chave Inglesa** | Repara a estação e a traz de volta ao funcionamento (+150 pts). |
| **Gato Travesso** *(Gato Laranja Roaming)* | Pula nos aparelhos e multiplica a queima em 2.5x. | **Brinquedo / Guizo** *(ou aproximação)* | Assusta o gato e o faz fugir para descansar (+50 pts). |

---

## 4. O Sistema "Panic & Recover"

### 4.1 A Barra de Caos (Chaos Meter)
* A barra no canto superior direito mede a tensão geral da cozinha (0% a 100%).
* **O que aumenta o Caos?**
  * Cada explosão que acontece adiciona instantaneamente **+25% de Caos**!
  * Aparelhos deixados em estado de Pânico (vermelho) vazam caos continuamente.
* **O que diminui o Caos?**
  * Desarmar aparelhos com sucesso reduz o Caos.
  * Realizar um **Clutch Recovery** alivia **-20% de Caos** de uma só vez!
* **Derrota (Game Over)**: Se a barra atingir **100%**, a cozinha sofre um *Kitchen Meltdown* (colapso total) e o restaurante é evacuado em chamas.

### 4.2 O que é "Clutch Recovery"?
Se você tiver a coragem de desarmar um aparelho nos **últimos 15% do seu tempo** (quando a barra estiver quase explodindo):
* Você realiza um **CLUTCH DEFUSE**!
* Ganha **+500 pontos** (em vez dos 100 normais).
* Dispara uma explosão de estrelas douradas e um som comemorativo triunfante.
* Reduz drasticamente o Caos da cozinha.

---

## 5. Pontuação e Avaliação de Desempenho

Ao término do turno de 2 minutos (ou em caso de colapso), o jogo exibe um relatório com:
* **Pontuação Final**
* **Total de Clutches Realizados**
* **Número de Explosões Ocorridas**
* **Recorde Pessoal (Top Chef Record)**
* **Classificação por Estrelas**:
  * ★★★ *(3 Estrelas)*: Sobreviver ao turno com mais de 3.000 pontos.
  * ★★☆ *(2 Estrelas)*: Mais de 1.500 pontos.
  * ★☆☆ *(1 Estrela)*: Pontuação inicial de participação.

---

## 6. Tabela Completa de Controles

O jogo suporta tanto **Teclado** quanto **Controle (Gamepad)** nativamente:

| Ação no Jogo | Teclado | Controle (Gamepad / Xbox / PS) |
| :--- | :--- | :--- |
| **Mover o Chef** | `W` / `A` / `S` / `D` ou `Setas` | Analógico Esquerdo ou D-Pad |
| **Interagir / Desarmar / Pegar / Espantar** | `Espaço` ou tecla `E` | Botão Sul (`A` / `X`) |
| **Largar Ferramenta Segurada** | `Q` ou `Shift` | Botão Oeste (`X` / `Quadrado`) |
| **Pausar o Jogo** | `Esc` ou `P` | Botão Start / Options |
| **Alternar Tela Cheia** | `F11` ou `Alt + Enter` | - |

---

## 7. Instruções de Execução

### Pré-requisitos
* Ter o **Go 1.26+** instalado na máquina.

### 7.1 Executar Diretamente no Computador (Desktop)
Abra o terminal na pasta do jogo e execute:
```bash
# Entrar na pasta do jogo
cd parte3-mini-game-jam/exploding-kitchens

# Executar o jogo
go run ./cmd/game
```

### 7.2 Rodar os Testes Unitários Automatizados
O projeto conta com suíte de testes cobrindo limites de movimentação do chef, máquina de estados dos eletrodomésticos, detecção de clutch e desalocação de partículas:
```bash
go test -v ./...
```

### 7.3 Compilar para WebAssembly (Navegadores)
O jogo é 100% compatível com a web e pode ser compilado para `.wasm`:
```bash
GOOS=js GOARCH=wasm go build -o game.wasm ./cmd/game
```

---

## 8. Destaques de Engenharia e Arquitetura (`ebitengineer`)

* **Zero Alocações no Loop de Desenho (`Draw`)**: Nenhum objeto ou fatia de memória é instanciado durante a renderização de quadros, garantindo **60 FPS sólidos** sem engasgos do Garbage Collector do Go.
* **Resolução Virtual 16:9 (`320x180`)**: Resolução fixa retrô que escala automaticamente sem distorções para janelas $1280 \times 720$ ou telas cheias em 4K.
* **Arte Procedural Pura (`procedural-art`)**: Todos os sprites do chef, eletrodomésticos, gato, ferramentas e efeitos visuais (fogo, fumaça, vapor, choque de explosão) são desenhados diretamente em código via buffers de pixels, eliminando dependência de imagens externas PNG/JPG.
* **Áudio Procedural em Código (`procedural-composer`)**: Síntese de áudio PCM de 16-bit estéreo a 44.1 kHz gerada em memória, sem arquivos WAV/MP3. A trilha sonora chiptune acelera o andamento (BPM) dinamicamente quando a cozinha entra em estado de pânico!
* **Modo Attract Demo (Estilo Fliperama)**: Se a tela de título ficar inativa por 10 segundos, uma IA assume o controle do chef no restaurante para demonstrar o jogo até que qualquer botão seja pressionado.
