# 📦 BoxBoy: Turbo Express (Edição Arcade Paperboy)

Um jogo de entrega e manobras em perspectiva pseudo-isométrica diagonal no estilo do clássico *Paperboy*, com customização avançada de personagem e veículo, sistema de Pânico & Recuperação com chefes urbanos e catástrofes de trânsito, e trilha sonora procedural com áudio DSP sintetizado em memória sem dependências externas!

---

## 🎮 Controles

| Ação | Teclado | Gamepad |
| :--- | :--- | :--- |
| **Manobra Lateral** | `A` / `D` ou `Setas Esquerda/Direita` | Analógico Esquerdo / D-Pad |
| **Acelerar / Frear** | `W` / `S` ou `Setas Cima/Baixo` | RT / LT |
| **Arremesso de Encomenda** | `Espaço` ou `J` ou Clique Esquerdo | Botão A (Sul) |
| **Bunny-Hop (Pular Guia/Buraco)**| `Shift` ou `K` | Botão X (Oeste) |
| **Buzina / Campainha** | `H`, `E` ou `B` | Botão B (Leste) |
| **Ação Heroica de Recuperação**| `R` ou `Enter` | Botão Y (Norte) |
| **Navegar Customização** | `W/S` (Categoria), `A/D` (Opção) | D-Pad |

---

## ✨ Funcionalidades Principais

1. **Perspectiva Isométrica Diagonal 26.5° (Paperboy)**:
   - A rota avança diagonalmente para cima e para a direita.
   - Pista central com asfalto e faixas, cercada por calçadas com casas residenciais e **Smart Lockers**.
   - Física balística 3D ($X, Y, Z$) para arremesso de encomendas sobre varandas, caixas de correio e lockers digitais.

2. **Customização Profunda do Entregador e Veículo**:
   - **Tons de Pele**: seleção visual por paleta de 5 cores (quadradinhos), sem rótulos textuais.
   - **Estilos de Cabelo**: Degradê Moderno, Coque Samurai, Black Power Afro, Moicano, Dreadlocks, Cabelo Longo.
   - **Cores de Cabelo**: Castanho Escuro, Loiro, Ruivo Fogo, Azul Neon, Rosa Choque Cyber, Prata.
   - **Uniformes Express**: Jaqueta Corta-Vento Amarela com Faixa Azul, Colete Refletivo Fluorescente, Moletom Streetwear Express, Polo Express.
   - **Acessórios de Cabeça**: Boné Aba Reta (virado pra trás), Viseira Retrô, Capacete Aerodinâmico, Bandana Ciclista.
   - **Óculos**: Óculos Ciclista Espelhado Neon, Óculos Nerd, Óculos Aviador Escuros.
   - **Mascotes do Cesto**: Cão Caramelo (mascote oficial), Capivara Zen, Mini Drone Rastreador.
   - **Veículos**: Bicicleta Urbana Amarela (Bunny-Hop ágil), Scooter Elétrica Express (Arrancada Turbo), Furgão Compacto Elétrico (Maior capacidade de carga).

3. **Ciclo de Chefes Urbanos: Pânico & Recuperação**:
   - **Boss 1 (km 0.6)**: *O Cliente Fantasma & A Matilha de Cães do Portão* (Pânico: matilha solta e aspersores; Recuperação: buzinar para afastar os cães e arremessar no Smart Locker).
   - **Boss 2 (km 1.3)**: *A Grande Panela de Asfalto (Cratera Voraz)* (Pânico: asfalto cedendo com erupções de lama; Recuperação: rampa de Bunny-Hop e estabilização com pacotes).
   - **Boss 3 (km 2.1)**: *O Megabloqueio da Manifestação Surpresa* (Pânico: barricadas e fumaça colorida; Recuperação: atalho na calçada e distribuição de brindes).
   - **Boss 4 (km 2.8)**: *O Tornado Metropolitano (SimCity Style)* (Pânico: vórtice de vento sugando pacotes; Recuperação: trava magnética de rodas e resgate de pacotes no ar).
   - **Boss 5 (km 3.4)**: *O Monstro do Atraso da Black Friday* (Pânico: esteiras transportadoras e relógio de contagem acelerado; Recuperação: Entrega Turbo contra os sensores do colosso).

4. **Áudio Procedural em Tempo Real (Zero Assets / Zero Disco)**:
   - Síntese pura DSP a 44.100 Hz estéreo de 16 bits.
   - Faixa BGM Principal: *Turbo Delivery Groove* (Funk-disco e synthwave brasileiro a 128 BPM com slap bass e metais).
   - Faixa BGM Pânico: *Panic Siren Theme* (152 BPM, sintetizadores tensos e sirene de alarme oscilante).
   - Fanfarra de Vitória: *Platinum Celebration Fanfare*.
   - Efeitos Sonoros: Arremesso *whoosh*, campainha residencial *ding-dong*, cascata de combos perfeitos, buzina dupla *triiim-triiim*, latidos, quebra de vidros e vento de tornado.

5. **HUD e Sistema de Reputação**:
   - Medidor de reputação com 5 Estrelas Douradas e selo *Rank Platinum*.
   - Contador de combo multiplicador de pontos.
   - Minimapa da rota em tempo real com marcadores dos 5 chefes.

---

## 🚀 Como Executar

```bash
# Compilar e rodar
go run .

# Ou compilar o executável binário autônomo
go build -o box-boy .
./box-boy
```

Para rodar os testes unitários:
```bash
go test -v ./...
```
