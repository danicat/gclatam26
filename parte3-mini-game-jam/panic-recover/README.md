# Panic Recover

Um arcade top-down sobre transformar um `panic()` em uma recuperação de emergência.

## Autoria

Eduardo Gutkoski

## Como jogar

1. No estado Calm, mova-se e escolha uma posição segura.
2. Pressione `Espaço` para ativar Panic, ou aguarde a contagem terminar.
3. Durante Panic, toque nos bugs para destruí-los.
4. Depois da meta de eliminações, alcance a zona `recover()` ciano antes que a estabilidade acabe.
5. Complete três ciclos para vencer.

Uma colisão durante Calm força um Panic com apenas 70% da estabilidade disponível.

## Controles

| Ação | Tecla |
| :--- | :--- |
| Mover | `WASD` ou setas |
| Ativar Panic / iniciar | `Espaço` ou `Enter` |
| Reiniciar | `R` |
| Sair | `Esc` |

## Executar

Requer Go 1.26.3 ou compatível:

```bash
go run .
```

Para validar o projeto:

```bash
go test ./...
go vet ./...
```

## Tecnologia e assets

- Go 1.26+ e Ebitengine v2.
- Canvas virtual 320×180, janela padrão 960×540.
- Sprite do gopher e sprite do bug gerados com Nano Banana.
- Música instrumental de 30 segundos gerada com Lyria 3 Clip.
- Efeitos sonoros sintetizados em Go e partículas com pool pré-alocado.

Os assets ficam embarcados no binário. A música é opcional: se não puder ser decodificada, o jogo segue com os efeitos procedurais.
