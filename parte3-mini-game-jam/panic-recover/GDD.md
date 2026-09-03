# Neon Breather — Game Design Document

> **Genre:** Arcade de arena top-down 2D  
> **Platform:** Desktop (Windows, macOS e Linux)  
> **Engine:** Go 1.26+ e Ebitengine v2  
> **Virtual resolution:** 320×180, proporção 16:9  
> **Target session:** 60–90 segundos  
> **Author / Lead Designer:** Eduardo Gutkoski

## 1. Visão geral

**Elevator pitch:** um gopher ativa Panic para ganhar velocidade e destruir bugs, mas precisa liberar e alcançar `recover()` antes que sua estabilidade acabe.

O tema combina o significado emocional de pânico e recuperação com referências literais a `panic()` e `recover()` de Go. A mecânica deve ser compreensível sem conhecimento de programação; mensagens e efeitos inspirados no runtime funcionam como uma segunda camada de humor para o público da GopherCon LATAM.

**Fantasia do jogador:** transformar um estado perigoso e parcialmente descontrolado em vantagem ofensiva, decidir quando assumir o risco e escapar no último instante.

**Referência de ritmo:** a inversão de poder de *Pac-Man* combinada com uma corrida curta de extração.

## 2. Pilares de design

1. **Panic é poder e perigo:** entrar em Panic torna o jogador ofensivo, mas inicia uma contagem irreversível.
2. **Recover precisa ser conquistado:** fugir não basta; o jogador deve destruir bugs antes de poder se recuperar.
3. **Caos legível:** velocidade, inércia, cores e áudio aumentam a tensão sem inverter comandos ou retirar a agência do jogador.
4. **Escopo de jam:** uma arena, um comportamento de inimigo e uma partida completa em menos de dois minutos.

## 3. Loop principal

Cada partida possui três ciclos de Calm → Panic → Recover:

1. O ciclo começa em **Calm**, com o jogador e todos os bugs daquele ciclo na arena.
2. Durante até 5 segundos, o jogador observa os bugs e escolhe uma posição.
3. `Espaço` ativa Panic antecipadamente. Caso contrário, Panic começa automaticamente ao fim da contagem.
4. Durante **Panic**, o jogador destrói bugs por contato enquanto sua estabilidade diminui por 12 segundos.
5. Ao atingir a meta de eliminações do ciclo, a zona `recover()` aparece em um ponto válido distante.
6. O jogador alcança a zona antes de a estabilidade chegar a zero.
7. Recover restaura a estabilidade e inicia o próximo ciclo após uma transição visual curta.
8. Alcançar Recover no terceiro ciclo vence a partida.

### Progressão dos ciclos

| Ciclo | Bugs na arena | Meta de eliminações | Mudança de dificuldade |
| :--- | ---: | ---: | :--- |
| 1 | 5 | 3 | Velocidade-base e perseguição suave |
| 2 | 8 | 5 | Mais bugs e velocidade maior |
| 3 | 11 | 7 | Maior quantidade e velocidade máxima |

Os números de velocidade serão constantes configuráveis e deverão ser ajustados por playtest. A meta, o número de bugs e os tempos acima são requisitos do MVP.

## 4. Mecânicas

### 4.1 Calm

- O jogador se movimenta normalmente, mas não destrói bugs.
- Bugs perseguem continuamente o jogador.
- Uma contagem visual de 5 segundos indica quando Panic será ativado automaticamente.
- Colidir com um bug força Panic imediatamente com 70% da estabilidade: 8,4 segundos.
- Ativação voluntária ou automática pelo cronômetro começa com 100% da estabilidade.

### 4.2 Panic

- O jogador destrói bugs ao tocá-los.
- A velocidade máxima aumenta.
- A desaceleração diminui progressivamente, produzindo derrapagem controlável.
- Os comandos permanecem consistentes e nunca são invertidos ou randomizados.
- A estabilidade esvazia linearmente durante 12 segundos, ou durante os 8,4 segundos restantes em um Panic forçado.
- Cada eliminação atualiza imediatamente a meta no HUD e produz partículas, flash e efeito sonoro.
- Chegar a estabilidade zero antes de tocar Recover encerra a partida em Game Over.

### 4.3 Recover

- A zona aparece somente depois que a meta de eliminações do ciclo é atingida.
- Ela surge no ponto válido mais distante do jogador entre oito posições candidatas fixas: quatro cantos e quatro pontos médios das bordas.
- O ponto deve respeitar uma margem das bordas e não pode se sobrepor ao jogador ou a um bug.
- Pulso ciano, feixe vertical e uma seta junto ao jogador indicam sua posição.
- Tocar a zona restaura a estabilidade e conclui o ciclo.
- Bugs restantes desaparecem durante a transição para evitar colisões injustas.

### 4.4 Bugs

- Existe um único comportamento: perseguir diretamente o jogador.
- Variações entre ciclos usam apenas velocidade, quantidade, tamanho e cor; não criam novas inteligências artificiais.
- Bugs não atacam à distância e não possuem vida. Uma colisão em Panic sempre os elimina.

## 5. Vitória, derrota e feedback

**Vitória:** concluir Recover nos três ciclos.

**Derrota:** a estabilidade chegar a zero durante Panic.

Não há vidas, pontuação, combo ou recorde no MVP. O resultado deve depender apenas do domínio do loop Panic → Recover.

Feedback obrigatório:

- Calm: arena e HUD verdes, música estável e movimento preciso.
- Panic: paleta vermelha, partículas, leve tremor de câmera e mensagens inspiradas em `panic: runtime error`.
- Recover disponível: destaque ciano, seta direcional e pulso sincronizado.
- Recover concluído: flash ciano e mensagem `recover successful`.
- Game Over: mensagem clara sobre estabilidade esgotada.

## 6. Controles

| Ação | Teclado |
| :--- | :--- |
| Mover | `WASD` ou setas |
| Ativar Panic | `Espaço` |
| Iniciar / confirmar | `Enter` ou `Espaço` |
| Reiniciar após o resultado | `R` |
| Sair | `Esc` |

Gamepad, mouse e toque ficam fora do MVP. As entradas devem ser representadas como ações lógicas para permitir extensões posteriores sem alterar as regras do jogo.

## 7. Direção visual e assets

**Estética:** terminal neon com fundo escuro, formas legíveis e humor de runtime.

**Cores funcionais:**

- Calm: verde.
- Panic: vermelho.
- Recover: ciano.
- Bugs: magenta ou âmbar, sempre contrastando com o estado atual.

**Nano Banana:** gerar um sprite estático do gopher e um sprite estático do bug, ambos adequados à escala virtual. Sprite sheets e animações quadro a quadro não fazem parte do MVP.

**Arte procedural em Go:** arena, grade de fundo, HUD, zona Recover, seta, flashes, glitches e partículas. Movimento, escala, rotação e modulação de cor dão vida aos sprites estáticos.

## 8. Áudio

**Música:** faixa eletrônica instrumental de aproximadamente 30 segundos criada com Lyria e reproduzida em loop. A composição deve sustentar tensão crescente sem depender de sincronização exata com os ciclos.

**Efeitos procedurais:**

- ativação voluntária de Panic;
- Panic forçado por colisão;
- eliminação de bug;
- liberação da zona Recover;
- Recover bem-sucedido;
- estabilidade crítica;
- Game Over e vitória;
- confirmação de interface.

Se a trilha do Lyria não estiver pronta a tempo, o jogo permanece completo com silêncio musical e efeitos procedurais. Ausência de música não pode bloquear a build.

## 9. Estados e interface

### Fluxo de cenas

`Title` → `Playing` → `Victory` ou `GameOver` → `Playing` por meio de `R`.

A tela inicial apresenta em três linhas:

1. `MOVE: WASD / ARROWS`
2. `PANIC: SPACE — DESTROY BUGS`
3. `REACH RECOVER BEFORE STABILITY HITS ZERO`

### HUD durante a partida

- Topo esquerdo: `CYCLE 1/3`.
- Topo central: barra de estabilidade.
- Topo direito: `BUGS 0/3`.
- Centro, temporariamente: contagem de Calm, aviso de Recover e resultados.
- Junto ao jogador: seta apontando para a zona Recover depois que ela for liberada.

## 10. Escopo técnico

- Arena fixa de uma tela, sem câmera de mundo ou tilemap.
- Atualização lógica a 60 TPS.
- Movimento independente da taxa de desenho.
- Colisões simples por círculos ou caixas delimitadoras, sem motor de física.
- Estado global da partida pequeno e explícito: cena, fase do ciclo, ciclo atual, estabilidade, contagem de eliminações, jogador, bugs e zona Recover.
- Sprites e áudio embarcados no binário com `go:embed` quando existirem.
- Assets obrigatórios devem falhar cedo com mensagem clara se estiverem inválidos.
- A lógica de regras e transições deve permanecer separada do desenho para permitir testes unitários.
- Não há backend, leaderboard, rede, persistência ou Cloud Run no MVP.

## 11. Validação e testes

Testes automatizados mínimos:

- ativação voluntária inicia Panic com 100% da estabilidade;
- colisão durante Calm inicia Panic com 70%;
- contagem de 5 segundos força Panic com 100%;
- metas de 3, 5 e 7 liberações são aplicadas aos ciclos corretos;
- Recover não aparece antes da meta;
- posição de Recover respeita margens e não sobrepõe entidades;
- estabilidade zero produz Game Over;
- terceiro Recover produz Victory;
- comandos nunca são invertidos durante Panic.

Playtest de aceitação:

- uma pessoa entende o objetivo apenas pela tela inicial;
- uma partida termina em 60–90 segundos;
- Calm, Panic e Recover são distinguíveis sem ler o HUD;
- a derrapagem aumenta a tensão sem parecer aleatória;
- um Panic forçado é difícil, mas ainda recuperável.

## 12. Fora do escopo

- chefes ou tipos adicionais de inimigo;
- fases, obstáculos, projéteis ou labirintos;
- vidas, pontuação, combos, ranking ou recordes;
- pausa, configurações ou seleção de dificuldade;
- gamepad, toque e build WebAssembly;
- backend ou serviços de nuvem;
- animações complexas ou sprite sheets;
- geração dinâmica de assets durante o jogo.

## 13. Ordem de corte

Se o tempo da jam estiver acabando, remover nesta ordem:

1. Música do Lyria, mantendo os efeitos procedurais.
2. Tremor, glitch e partículas secundárias.
3. Sprites gerados, substituindo-os por formas procedurais.
4. Tela de vitória dedicada, mantendo uma mensagem sobre a arena.

O loop Calm → Panic → Recover, os três ciclos, os controles e as condições de vitória e derrota não podem ser cortados.
