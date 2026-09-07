package entities_test

import (
	"testing"

	"box-boy/internal/entities"
)

func TestBossListAndPanicRecover(t *testing.T) {
	bosses := entities.NewBossList()

	if len(bosses) != 5 {
		t.Fatalf("esperava 5 bosses temáticos de pânico e recuperação, obteve %d", len(bosses))
	}

	b1 := bosses[0]
	if b1.Type != entities.BossDogPack {
		t.Errorf("primeiro boss deveria ser BossDogPack, obteve %v", b1.Type)
	}

	// Simula aproximação do jogador
	b1.Update(0.1, b1.TriggerY-50.0)
	if b1.State != entities.BossApproaching {
		t.Errorf("esperava estado BossApproaching, obteve %v", b1.State)
	}

	// Jogador atinge o gatilho: PÂNICO!
	b1.Update(0.1, b1.TriggerY+10.0)
	if b1.State != entities.BossPanic {
		t.Errorf("esperava estado BossPanic ao atingir o marco, obteve %v", b1.State)
	}
	if b1.ScreenShake <= 0 {
		t.Errorf("tremor de tela (ScreenShake) deveria estar ativo no pânico")
	}

	// Ações heroicas de recuperação
	for i := 0; i < b1.RecoverTarget-1; i++ {
		defeated := b1.AddRecoverAction()
		if defeated {
			t.Errorf("não deveria derrotar o chefe antes de atingir todas as ações de recuperação")
		}
	}

	// Última ação de recuperação: Derrota triunfante do chefe!
	defeated := b1.AddRecoverAction()
	if !defeated {
		t.Errorf("deveria ter derrotado o chefe após completar a meta de recuperação")
	}
	if b1.State != entities.BossDefeated {
		t.Errorf("esperava estado BossDefeated, obteve %v", b1.State)
	}
}
