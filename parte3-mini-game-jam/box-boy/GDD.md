# Game Design Document (GDD)
# BoxBoy: Turbo Express (Entrega Turbo)

> **Game Title**: BoxBoy: Turbo Express  
> **Target Genre**: Arcade Delivery Action / Pseudo-Isometric Runner (Estilo Paperboy)  
> **Target Platform**: Desktop (Windows/macOS/Linux) & WebAssembly (Google Cloud Run / Browser)  
> **Target Aspect Ratio & Resolution**: 16:9 Widescreen (`640x360` Canvas Virtual Pixel Scaling)  
> **Target Technologies**: Go 1.26+, Ebitengine v2 (`github.com/hajimehoshi/ebitengine/v2@latest`), Antigravity (`agy`), Nano Banana (Sprites/Pixel Art), Lyria 3 (BGM)  
> **Author / Lead Designer**: Antigravity Game Designer Agent  

---

## 1. Visão Geral & Pitch Executivo

### 1.1 Elevator Pitch
**BoxBoy: Turbo Express** é uma reimaginação arcade no estilo clássico de *Paperboy* com estética vibrante e energia urbana. O jogador pilota veículos de entrega personalizáveis (bike, scooter elétrica ou van compacta), arremessando encomendas de caixas amarelas diretamente nas portas, caixas de correio e lockers inteligentes de clientes ao longo de bairros movimentados. Durante o percurso, o jogador enfrenta imprevistos urbanos e catástrofes ao estilo *SimCity*, enfrentando Bosses baseados no ciclo de **Pânico & Recuperação** (crises caóticas que exigem manobras heroicas e triagens ágeis para salvar a entrega e restaurar a ordem).

### 1.2 Inspirações Centrais
- **Paperboy (Arcade/NES)**: Movimentação diagonal/isométrica contínua, arremesso com timing e pontaria em caixas de entrega, obstáculos dinâmicos na calçada e na rua.
- **SimCity (Disasters)**: Catástrofes aleatórias e bosses de forças da natureza ou do caos urbano (tornados metropolitanos, crateras de asfalto, manifestações).
- **Crazy Taxi & Micro Machines**: Sensação de velocidade, desvios milimétricos, buzinas e física satisfatória de veículos leves.
- **Cultura Express & Delivery Latino-Americano**: Uniforme amarelo canário, fita gomada azul, carros elétricos ágeis, cachorros caramelo no portão, interfones e a promessa de "Chega Amanhã (ou em 5 minutos!)".

### 1.3 Audiência e Clima (Mood)
- **Tom**: Divertido, acelerado, bem-humorado, caótico e satisfatório.
- **Ritmo**: Sessões rápidas de 3 a 5 minutos por rota/fase, com picos intensos de adrenalina durante as fases de Pânico dos Bosses e momentos catárticos de Recuperação.

---

## 2. Loop de Gameplay Principal & Mecânicas

### 2.1 O Loop Principal (Action $\rightarrow$ Challenge $\rightarrow$ Reward)
```text
[Seleção de Veículo & Carga]
          │
          ▼
[Navegar Bairro em Rolagem Diagonal] ──(Obstáculos: Pedestres, Buracos, Pets)
          │
          ├─► [Mirar & Arremessar Encomenda no Alvo] ──► (Porta, Varanda, Smart Locker)
          │         │
          ├─ Sucesso Perfeito: +Pontos Turbo, Reputação Ouro, Moedas
          └─ Erro (Janela Quebrada / Piscina): -Reputação, Pacote Danificado
          │
          ▼
[EVENTO DE BOSS: CICLO DE PÂNICO & RECUPERAÇÃO]
  ├─ Fase Pânico: O ambiente colapsa, sirenes soam, caos total e perda de controle iminente.
  └─ Fase Recuperação: Ações de salvamento (desvios precisos, triagem de pacotes, reparos turbo).
          │
          ▼
[Chegada ao Ponto Final / CDD de Bairro] ──► [Avaliação 5 Estrelas & Customização]
```

### 2.2 Mecânicas de Condução e Arremesso
1. **Pilotagem Dinâmica**:
   - Controle de velocidade (Acelerar, Frear, Drift leve na curva).
   - Manobras evasivas: Salto/Bunny-Hop (para bike/scooter pular guias e buracos) e Buzina (afasta pedestres e pets).
2. **Sistema de Arremesso de Pacotes**:
   - Mira com indicador de trajetória baseado na velocidade do veículo.
   - 3 Tipos de Encomendas:
     - **Pacote Padrão (Caixa Pequena)**: Leve, voa rápido, ideal para caixas de correio.
     - **Pacote Frágil (Eletrônico com fita de cuidado)**: Exige arremesso suave em tapetes de entrada; estoura se bater em paredes.
     - **Pacote Pesado / Especial**: Arremesso curto com parábola, precisa cair em garagens ou lockers abertos.
3. **Mecânica de Reputação do Entregador**:
   - Barra de Reputação (Termômetro de Qualidade):
     - **Nível Ouro/Platina**: Entrega perfeita no alvo sem danos.
     - **Nível Alerta**: Pacotes atrasados ou danificados reduzem a pontuação e atraem clientes irritados.

---

## 3. Sistema Rico de Customização do Entregador & Veículos

Antes de iniciar a jornada (e acessível na garagem entre turnos), o jogador conta com um criador de personagem profundo em pixel art:

```text
┌────────────────────────────────────────────────────────────────────────┐
│                        CENTRAL DO ENTREGADOR                           │
├──────────────────────────────┬─────────────────────────────────────────┤
│         PREVIEW 2D           │           CATEGORIAS DE ITENS           │
│                              │                                         │
│            (o_o)             │  [1] Tom de Pele & Rosto                │
│            /|M|\  <-- Uniforme│  [2] Cabelo & Barba (Cores Neon/Retro) │
│             / \   <-- Calça  │  [3] Bonés, Capacetes & Viseiras        │
│          =========           │  [4] Jaquetas, Coletes & Moletons       │
│         [VEÍCULO 2D]         │  [5] Mochila / Baú Térmico Custom       │
│        (O)=====(O)           │  [6] Acessórios (Óculos, Tatuagem, Pins)│
│                              │  [7] Mascote do Cesto (Caramelo, Drone) │
│                              │  [8] Veículo (Bike, Scooter, Van Turbo) │
└──────────────────────────────┴─────────────────────────────────────────┘
```

### 3.1 Opções de Customização do Entregador
- **Biótipo & Pele**: Tons de pele variados, formato de rosto, expressões faciais (determinado, sorridente, óculos escuros estilosos).
- **Cabelo & Barba**: Dreadlocks, corte degradê, crespo volumoso, coque samurai, moicano punk, liso, calvície charmosa, cores vibrantes.
- **Vestuário do Entregador**:
  - Jaqueta Corta-Vento Amarela com faixas refletivas.
  - Colete Utilitário de Entrega Rápida com bolsos e crachá personalizado.
  - Moletom Streetwear com capuz e logo estilizado.
  - Bermudas e calças táticas com joelheiras de proteção.
- **Acessórios & Estilo**:
  - Boné virado para trás, capacete aerodinâmico esportivo ou viseira retrô.
  - Óculos de sol estilo ciclista (espelhado neon), óculos de grau nerd-cool.
  - Tênis de corrida com molas ou botas reforçadas.
  - Luvas de ciclista sem dedos com faixas fluorescentes.
- **Mascote da Garupa / Cesto Dianteiro**:
  - *Caramelo*: O lendário vira-lata brasileiro que late e fareja atalhos.
  - *Capivara Zen*: Fica calma no cesto, reduzindo o tempo da fase de Pânico.
  - *Mini Drone de Apoio*: Ilumina a rua e avisa sobre buracos à frente.

### 3.2 Veículos Disponíveis
| Veículo | Agilidade | Velocidade Máx | Capacidade de Carga | Habilidade Especial |
| :--- | :--- | :--- | :--- | :--- |
| **Bicicleta Urbana Amarela** | ★★★★★ | ★★☆☆☆ | 10 Pacotes | Salto de Guia & Passagem por Calçadas Estreitas |
| **Scooter Elétrica Express** | ★★★★☆ | ★★★★☆ | 16 Pacotes | Turbo Boost de Bateria & Buzina Ultrassônica |
| **Furgão Compacto / Buggy Elétrico** | ★★☆☆☆ | ★★★★★ | 35 Pacotes | Aríete contra pequenos obstáculos & Lançador Triplo |

---

## 4. O Sistema de Bosses: Abstrações de "Pânico & Recuperação"

Em vez de usar jargões de código ou conceitos artificiais, o jogo traduz o paradigma de **Pânico** (desastre imprevisto que quebra a rotina) e **Recuperação** (sangue frio e engenhosidade humana para contornar o problema) em eventos monumentais inspirados em perrengues reais de entregadores e catástrofes de *SimCity*.

```text
        ┌─────────────────────────────────────────────────────────┐
        │             MECÂNICA DE PÂNICO & RECUPERAÇÃO            │
        └─────────────────────────────────────────────────────────┘
                                     │
           [NORMAL]: Entrega ritmada, música groove fluida
                                     │
                                     ▼  (Gatilho do Boss)
           ┌─────────────────────────────────────────────────────┐
           │                     FASE PÂNICO                     │
           │  • Sirenes e distorção no áudio                     │
           │  • Visão periférica treme, chuva de detritos         │
           │  • Controles ficam mais sensíveis ou escorregadios   │
           │  • Rota principal destruída ou bloqueada            │
           │  • Carga de pacotes corre risco de ser perdida      │
           └─────────────────────────────────────────────────────┘
                                     │
                    (Ação de Resposta do Jogador)
                                     ▼
           ┌─────────────────────────────────────────────────────┐
           │                  FASE RECUPERAÇÃO                   │
           │  • Respiração profunda / Foco turbo ativado         │
           │  • Minigame de alinhamento / desvios em sequência   │
           │  • Resgate acrobático dos pacotes que voaram no ar  │
           │  • Uso de atalhos e itens de suporte                │
           │  • Vitória restaura o bairro e concede Reputação 5★ │
           └─────────────────────────────────────────────────────┘
```

### 4.1 Lista de Bosses & Mini-Bosses

#### 1. Mini-Boss: O "Cliente Fantasma & A Matilha do Portão"
- **Contexto**: A campainha toca, ninguém atende, o aplicativo acusa "destinatário ausente", e o portão automático emperra soltando 5 cachorros bravos e canhões de aspersores de jardim!
- **Fase Pânico**:
  - Latidos ensurdecedores em surround, água espirrando na lente da câmera embaçando a visão, pacotes escorregando do bagageiro.
- **Fase Recuperação**:
  - O jogador arremessa petiscos de brinquedo para acalmar a matilha, executa um salto perfeito sobre os jatos de água e joga o pacote exatamente no Smart Locker seguro antes do tempo esgotar!

#### 2. Mini-Boss: A Grande Panela de Asfalto (A Cratera Voraz)
- **Contexto**: Uma obra malfeita na avenida cede de repente, abrindo uma fenda vulcânica de lama, tubulações estouradas e placas de sinalização voadoras.
- **Fase Pânico**:
  - O asfalto se quebra em blocos flutuantes, bueiros disparam colunas de água pressurizada que desviam a trajetória das encomendas.
- **Fase Recuperação**:
  - Manobra de drift precisa pelas beiradas de concreto, uso de rampas de entulho para saltar sobre a fenda, jogando sacos de brita para estabilizar o solo e salvar a carga intacta.

#### 3. Boss Maior: O Megabloqueio da Manifestação Surpresa
- **Contexto**: Uma multidão caótica, caminhões de som com buzinas ensurdecedoras, bandeiras gigantescas e labirintos de barricadas fecham 4 quarteirões inteiros.
- **Fase Pânico**:
  - Ruído ensurdecedor de vuvuzelas e apitos; fumaça colorida reduz o campo de visão; pedestres empurram carrinhos de compras na direção do jogador.
- **Fase Recuperação**:
  - O entregador encontra passagens secretas por dentro de galerias comerciais e becos floridos, faz entregas rápidas de água gelada e brindes para abrir passagem pacífica no meio da multidão, emergindo triunfante no outro lado da avenida!

#### 4. Boss Maior: O Tornado Urbano Metropolitano (Estilo SimCity)
- **Contexto**: Uma tempestade subtropical repentina forma um vórtice gigante de vento que desce pela avenida comercial, sugando toldos, cadeiras de plástico e placas de trânsito.
- **Fase Pânico**:
  - Força centrífuga arrasta o veículo para o olho do tornado; caixas de papelão começam a ser sugadas para o céu; o vento altera a trajetória física dos arremessos em 90 graus!
- **Fase Recuperação**:
  - Ativar a trava eletromagnética das rodas da bike/van no asfalto; arremessar os pacotes contra o vento usando o ricochete de prédios; capturar no ar as encomendas desgarradas com manobras aéreas estilizadas!

#### 5. Mega Boss Final: O Monstro do Atraso da "Black Friday"
- **Contexto**: Uma personificação colossal do caos logístico — um monstro mecânico formado por montanhas de encomendas acumuladas, esteiras rolantes rebeldes, relógios de ponto girando furiosamente e raios de tempestade.
- **Fase Pânico**:
  - A contagem regressiva para a "Meia-Noite" pisca na tela; o monstro bate esteiras no chão criando ondas de choque; encomendas voam em todas as direções prontas para quebrar.
- **Fase Recuperação**:
  - O jogador aciona o "Modo Entrega Turbo Relâmpago", disparando caixas nos pontos fracos (dutos de leitura de código de barras do monstro), neutralizando as esteiras com saltos milimétricos e entregando a última encomenda de ouro diretamente na mão da cliente idosa no topo do arranha-céu!

---

## 5. Estratégia de Áudio & Soundscape ("Entregar Tudo")

O áudio deve ser uma das identidades mais marcantes do jogo: contagiante, dinâmico e reagindo instantaneamente a cada sucesso, arremesso e crise.

### 5.1 Trilha Sonora Musical (BGM - Gerada via Lyria 3)
- **Gênero**: Funk Carioca / Nu-Disco / Synthwave Híbrido em 128 BPM com metais enérgicos, slap bass pulsante e batidas tropicais contagiantes.
- **Camadas Interativas Dinâmicas**:
  1. *Camada Base (Passeio Urbano)*: Groovy suave, percussão leve com violão e sintetizador ensolarado.
  2. *Camada Ação (Entregando)*: Bumbo potente entra em cena, trompetes celebram combos de entregas perfeitas.
  3. *Camada Pânico (Boss/Catástrofe)*: O ritmo acelera drasticamente para 150 BPM, sintetizadores distorcidos, sirenes estilizadas em clave musical, baixo pesado e sensação de urgência máxima.
  4. *Camada Recuperação (Virada Heroica)*: Solos triunfantes de saxofone/trompete, retorno triunfal do ritmo principal com coro de comemoração!

### 5.2 Efeitos Sonoros (SFX - Síntese Procedural Pure-Code & Amostras Cristalinas)
- **Arremesso de Caixa**: *Whoosh* aerodinâmico suave seguido pelo som satisfatório de papelão acolchoado.
- **Acerto Perfeito no Alvo**: Chime cristalino tipo caixa registradora de entrega confirmada + *ding-dong* harmonioso de campainha.
- **Erro / Vidro Quebrado**: Efeito cômico de *crash*, espirro d'água na piscina com patinho de borracha *quack*.
- **Buzina do Veículo**: Triim-triiim clássico de campainha de bike ou bi-bi estridente da motinho elétrica.
- **Pânico Sfx**: Alarme de proximidade pulsante, freada com pneus cantando no asfalto, vento uivante de tornado.
- **Vitória de Rota**: Jingle icônico de notificação de "Pacote Entregue com Sucesso" em versão orquestrada 16-bit!

---

## 6. Esquema de Controles & Mapeamento de Entrada

Suporte nativo e balanceado para Teclado, Gamepad e Touchscreen (WebAssembly Mobile/Tablet):

| Ação Lógica | Teclado & Mouse | Controle (Xbox / DualShock) | Touchscreen (WASM Mobile) |
| :--- | :--- | :--- | :--- |
| **Mover / Direcionar** | `W`/`S` ou `↑`/`↓` (Pista) | D-Pad / Analógico Esquerdo | D-Pad Virtual na esquerda |
| **Acelerar / Frear** | `D` (Acelerar) / `A` (Frear) | Gatilhos `RT` / `LT` | Botões Virtuais de Aceleração |
| **Arremessar Encomenda** | `Barra de Espaço` / Clique Esq. | Botão `A` (Sul) / `RB` | Botão Flutuante de Arremesso |
| **Salto / Pulo de Guia** | `W` + Espaço ou `Shift` | Botão `X` (Oeste) | Botão "Pular Guia" |
| **Buzina / Afastar Pets** | `H` ou `E` | Botão `B` (Leste) | Botão de Buzina |
| **Ação de Recuperação** | `R` ou `Enter` (QTE / Foco) | Botão `Y` (Norte) | Ícone de Ação Heroica |
| **Pausa / Menu** | `Esc` / `P` | Botão `Start` / `Options` | Ícone de Engrenagem (Topo Direito) |

---

## 7. Interface de Usuário (HUD) & Fluxo de Telas

### 7.1 Layout da HUD Durante o Jogo
```text
┌────────────────────────────────────────────────────────────────────────┐
│ [5★] Reputação: 98% [Ouro]        📦 Carga: 12/20        ⏱️ Tempo: 02:45 │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│                                              [Casa com Pedido]         │
│                                                Alvo: Varanda           │
│                 (Veículo / Entregador)       ┌──────────────┐          │
│                      \    O===O              │ 🎯 [DROP]    │          │
│                       \--[ 📦 ]              └──────────────┘          │
│                                                                        │
├────────────────────────────────────────────────────────────────────────┤
│ [MINI-MAPA DE ROTA: ●───●───▲(Boss)───🏁]       [COMBO: 4x PERFEITO!]  │
└────────────────────────────────────────────────────────────────────────┘
```
- **Topo Esquerdo**: Nível de Reputação (Estrelas de 1 a 5) e feedback de satisfação do cliente.
- **Topo Centro**: Quantidade de encomendas restantes no baú / mochila.
- **Topo Direito**: Cronômetro de entrega expressa.
- **Rodapé Esquerdo**: Barra de progresso linear da rua (mostrando pontos de entrega e ícone de alerta de Boss).
- **Rodapé Direito**: Contador de Combo de entregas perfeitas consecutivas.

### 7.2 Fluxo de Telas (Scene Flow)
```text
[Splash Logo: Antigravity & BoxBoy]
               │
               ▼
   [Menu Principal / Modo Attract Demo]
               │
               ├─► [Central de Customização do Entregador] (Salva perfil local)
               │
               ├─► [Garagem de Veículos] (Upgrade e seleção de Bike/Scooter/Van)
               │
               ▼
       [Gameplay: Rua da Entrega]
               │
        (Gatilho de Boss)
               ▼
     [Evento Pânico / Recuperação]
               │
               ├─► Sucesso ──► [Conclusão da Rota & Avaliação 5 Estrelas]
               └─► Falha   ──► [Game Over: Reputação Zerada & Re-tentativa Rápida]
```

---

## 8. Arquitetura Técnica & Especificação Ebitengine v2 (Go 1.26+)

### 8.1 Estrutura Modular de Pacotes (`internal/`)
```text
box-boy/
├── cmd/
│   └── game/
│       └── main.go               # Ponto de entrada Ebitengine (SetWindowSize, RunGame)
├── internal/
│   ├── assets/                   # Embed FS para sprites, fontes TTF e áudio
│   │   ├── embedded.go
│   │   └── data/
│   ├── config/                   # Resolução virtual 640x360, constantes de física e tuning
│   ├── customizer/               # Sistema de personalização de personagem e paletas
│   ├── entities/                 # Entidades: Player, Vehicle, Package, House, Obstacle, Boss
│   │   ├── player.go
│   │   ├── vehicle.go
│   │   ├── package.go
│   │   ├── house.go
│   │   ├── obstacles.go
│   │   └── bosses.go             # Implementação dos chefes e estados Panic/Recover
│   ├── input/                    # Mapeamento abstrato de teclado, gamepad e touch
│   ├── physics/                  # Colisões AABB e Spatial Hash de alta performance
│   ├── render/                   # Renderização pseudo-isométrica com ordenação por profundidade Y
│   ├── scenes/                   # FSM de Cenas: Title, Customizer, GamePlay, GameOver, Victory
│   ├── audio/                    # Player de BGM e sintetizador DSP de SFX procedurais
│   └── ui/                       # HUD flexível, barras de reputação, minimapa e fontes
├── GDD.md                        # Este documento
├── go.mod
└── go.sum
```

### 8.2 Diretrizes Técnicas Ebitengine & Projeção Isométrica
1. **Projeção Isométrica Diagonal (Paperboy Style)**:
   - A rua avança em vetor diagonal ascendente (ângulo clássico 26.5° / razão 2:1).
   - Matriz de transformação de coordenadas mundo -> tela:
     $$X_{tela} = (X_{mundo} - Y_{mundo}) \cdot \cos(26.5^\circ) + CamX$$
     $$Y_{tela} = (X_{mundo} + Y_{mundo}) \cdot \sin(26.5^\circ) - Z_{altura} + CamY$$
   - Ordenação de profundidade (Y-Sorting): Todas as entidades (casas, postes, entregador, pacotes no ar, bosses) são desenhadas ordenadas por $Y_{tela} + Z_{base}$ para profundidade correta sem artefatos visuais.
2. **Resolução Virtual Fixa (16:9)**:
   - `virtualWidth = 640`, `virtualHeight = 360` (garante nitidez perfeita para pixel art e visibilidade ampla da rua e das calçadas laterais).
   - `Layout()` garante redimensionamento perfeito em qualquer monitor e telas mobile sem distorção.
3. **Ciclo de Tempo Delta ($dt$) e 60 FPS**:
   - Movimentação de veículos e projéteis de pacotes calculados com delta time real para manter fidelidade física mesmo se houver oscilação de quadros no navegador.
4. **Zero-Allocation no Draw Loop**:
   - Reutilização estrita de `ebiten.DrawImageOptions` e `ebiten.GeoM` com buffers pré-alocados para evitar pressão no Garbage Collector do Go.
5. **Estratégia de Assets em Duas Fases**:
   - **Fase 1 (Protótipo Imediato Zero-Dependency)**: Implementação de todo o gameplay, física de arremesso, casas, veículo inicial (Bicicleta Amarela com Bunny-Hop) e partículas com `procedural-art` (Go puro gerando texturas em memória) e `procedural-composer` para SFX.
   - **Fase 2 (Polimento de Produção)**: Substituição e enriquecimento com spritesheets gerados via `nano-banana` (animações completas de pedalada, customização e bosses) e trilhas sonoras orquestradas em alta definição via `lyria 3`.
