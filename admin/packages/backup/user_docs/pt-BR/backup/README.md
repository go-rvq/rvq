# Backup

Uma cópia da **aplicação inteira** — o banco de dados, os arquivos do site,
os uploads, o próprio programa, a sua configuração —, cada uma um **state**.
Um cliente baixa o state num diretório, e a aplicação sobe nele.

Um state é feito no próprio servidor, sem parar a aplicação: os arquivos
são ligados (hard links) e não copiados — um state ocupa no disco só o que
mudou desde então —, e um **manifesto** lista cada arquivo com o seu sha256.
O cliente compara o manifesto com o que já tem e baixa só o que mudou.

## Um state novo

**Novo** pede:

- **Descrição**: por que este state ("antes da migração 2.0").
- **Fontes**: o que entra. Sem alguma, o state é **parcial**: serve para
  recuperar uma parte, mas não sobe a aplicação sozinho.
- **Manter**: a retenção nunca o apaga.

Ele entra **na fila** e é **gerado** em segundo plano; a lista e o detalhe
mostram o progresso de cada fonte. Fica **pronto**, ou **falhou** (o erro no
detalhe), ou foi **cancelado**.

Um state também é feito sozinho, todo dia (origem **agendado**). Dos
states prontos, ficam os 7 mais novos (e outros tantos dos que falharam ou
foram cancelados); os marcados **Manter** ficam sempre.

## O detalhe

- **Estado**, **Origem**, **Criado por**, **Início** e **Fim**, **Arquivos** e
  **Tamanho**, **Manifesto (sha256)**.
- **Por fonte**: os arquivos, o tamanho, o tempo e os avisos de cada fonte.
- **Download**: pronto, o link do manifesto e as **rotas do cliente**, com o
  id do state.
- **Downloads**: as sessões de download abertas para ele.

Um state sendo gerado não pode ser excluído: cancele-o antes. Um state cujos
arquivos sumiram (a pasta removida no servidor) é marcado como falho ao
reiniciar a aplicação.

## Os arquivos de um state

**Arquivos**, no detalhe de um state pronto, abre os seus arquivos no
explorador, só para leitura: navegue pelas pastas, veja um texto ou uma
imagem e baixe.

- O ícone de download da barra baixa a pasta aberta em `.tar.gz`; com
  itens selecionados, o que estiver selecionado.
- Um arquivo vem como ele é; uma pasta, num `.tar.gz`.
- Mais de um item: um menu deixa escolher entre **um `.tar.gz`** com todos
  e **um a um**.
- O menu de contexto (clique direito) faz o mesmo: na área vazia, a pasta
  aberta; num item ou na seleção, o que estiver selecionado.

Quem pode ver o state pode ler todos os seus arquivos.

## Um link para baixar tudo

**Link de download (.tar.gz)**, no detalhe, cria um link de 24 horas que
baixa o state inteiro sem login (veja a ação). Os links criados ficam em
**Links**, onde podem ser revogados.

## Baixar um state

Um programa baixa um state com uma **chave de acesso** (o seu perfil ›
**Chaves de acesso**) que não negue a permissão **Baixar** desta página.
A chave abre uma **sessão**, que dá um **token**; o token baixa o manifesto
e os arquivos — e renova a sessão — por uma hora.

```
ADMIN=https://seu-site/admin
KEY=hck_…                                   # a chave de acesso

# o state pronto e completo mais novo
curl -fsS -H "Authorization: Bearer $KEY" "$ADMIN/!backup/latest"
# → {"id":"…","created":"…","files":…,"bytes":…}

# uma sessão: o token
curl -fsS -X POST -H "Authorization: Bearer $KEY" "$ADMIN/!backup/<id>/session"
# → {"token":"…","expires":"…","state":"<id>","manifest":"…","files":"…"}

# o manifesto e cada arquivo, pelo token
curl -fsS -H "Authorization: Bearer $TOKEN" "$ADMIN/!backup/<id>/manifest" -o manifest.gz
curl -fsS -H "Authorization: Bearer $TOKEN" "$ADMIN/!backup/<id>/files/<caminho>" -o <caminho>
```

Um arquivo aceita **Range** (um download interrompido continua de onde
parou) e **If-Match** com o seu sha256 (o ETag): se o arquivo mudou, a
resposta é 412 em vez de um arquivo misturado.

O manifesto é um texto compactado (gzip), uma linha por entrada, separada
por tabulações: o sha256 (base64url, sem `=`), o modo, a data, o tamanho, o
tipo (`f` arquivo, `d` pasta, `l` link) e o caminho.

## Permissões

- **Ver** a página: a lista, o detalhe, a página **Arquivos** (ler e baixar os arquivos de um state), o manifesto pelo admin.
- **Criar**, **Editar** (a descrição, **Manter**), **Excluir**, **Cancelar**.
- **Baixar** (`!download`): o cliente (o manifesto e os arquivos por uma
  chave de acesso) e o **Link de download** — todos os dados e segredos da
  aplicação. Dê-a com cuidado.

## Caminhos estáveis

O cliente usa o serviço `!backup`, fora das páginas do menu e do login do
admin (veja o guia **Caminhos estáveis**). Os endereços não mudam com o
menu. `<admin>` é o endereço do admin.

| Método | Caminho | Acesso | O que dá |
| --- | --- | --- | --- |
| GET | `<admin>/!backup/latest` | chave de acesso | o id do state pronto e completo mais novo |
| POST | `<admin>/!backup/{id}/session` | chave de acesso, ou o token (renova) | uma sessão: o seu token, os endereços do manifesto e dos arquivos |
| GET | `<admin>/!backup/{id}/manifest` | token, ou a sessão do admin | o manifesto |
| GET | `<admin>/!backup/{id}/files/{caminho}` | token | um arquivo do manifesto (Range, If-Match) |
| GET | `<admin>/!backup/{id}/archive/{token}.tar.gz` | link | o state inteiro em `.tar.gz` (24 horas; sem Range) |
| GET | `<admin>/!backup/{id}/fs/api/fs/{op}?path=…` | sessão | o explorador da página Arquivos: listar, ler, baixar (só leitura) |

O link de download não pede `Authorization`: o próprio link é o acesso. Nos outros caminhos, sem o cabeçalho `Authorization`, um pedido vai ao admin: pela sessão do
admin, só o manifesto (o link do detalhe). Respostas de erro são JSON
(`{"error": "…"}`): 401 sem chave ou com uma recusada, 403 se a chave não
pode baixar, 404 se o state não existe, 409 se ainda não está pronto.
