# publish_test

Testes do exemplo de publish (`examples_admin.PublishExample`) dirigidos pelos
eventos que os botões do admin disparam, conferidos pelo que ficam no banco e
pelo que a resposta manda a página fazer — não pelo texto exato do JavaScript.

| Teste | Cenário |
|---|---|
| `TestNew` | criar: primeira versão, sem pai, rascunho |
| `TestDuplicate` | duplicar (inclusive uma versão publicada) e duplicar de novo |
| `TestPublishAndUnpublish` | publicar; publicar outra versão tira a anterior do ar; despublicar |
| `TestSchedule` | agendar (`agora < início < fim`) e as recusas, que mantêm o agendamento salvo |
| `TestVersionDialog` | renomear, listar as versões do registro, aba de nomeadas, busca, escolher uma versão |
| `TestDeleteVersion` | apagar outra versão, a exibida (passa à mais antiga), a mais antiga exibida (passa à mais nova) e a última (volta à listagem) |

Substituem os fluxos gravados no browser pelo gerador do qor5
(`testflow/gentool`), que conferiam strings de JS que o presets não emite mais.
