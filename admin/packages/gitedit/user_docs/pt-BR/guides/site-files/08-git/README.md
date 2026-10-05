# Pelo git

Os mesmos arquivos podem ser clonados e enviados (push) com o **git**, do seu
computador, com o seu **login e senha** do admin. Há dois repositórios — as
URLs estão na página {%= admin.page("/site-files").link %}, em **Pelo git**:

| Repositório | Um push… | Pede |
| --- | --- | --- |
| **o seu rascunho** — `…/site-files-draft.git` | muda o rascunho: o editor e o **preview** mostram na hora; depois **publique** pela página | `!git` |
| **o site** — `…/site-files.git` | **publica** direto | `!git` e `!publish` |

```
git clone https://seu-site/admin/site-files-draft.git
cd site-files-draft
# … edite, commit …
git push
```

## O que um push confere

Cada arquivo alterado pede **o que o editor pede** — as mesmas permissões, de
cada caminho (veja o passo **Permissões**):

| No push | Pede |
| --- | --- |
| um arquivo novo | `!create` |
| um arquivo alterado | `!edit` |
| um arquivo excluído | `!delete` |
| renomeado na mesma pasta | `!rename` (na origem e no destino) |
| movido para outra pasta | `!move` (na origem e no destino) |

Uma negação recusa o push inteiro — nada é gravado — e o git mostra os
arquivos: `remote: config/layout_config.gad: no permission (!edit)`.

Além disso:

- **Só a branch do site**, e só para frente (fast-forward): se outros
  publicaram, faça `git pull` antes.
- **No rascunho**, o push é recusado enquanto houver alterações sem commit
  feitas no editor: faça o commit ou descarte-as antes.
- **No site**, o push confere o mesmo que **Publicar**: os templates compilam,
  e os arquivos do site não foram mudados fora do git. Um push e uma
  publicação pela página nunca acontecem ao mesmo tempo.

## Quem fez cada commit

Cada commit enviado precisa dizer **quem o fez neste site**: uma linha no fim
da mensagem, como a `Co-authored-by` do git —

    Site-User: o-seu-login <a-sua-chave@o-endereço-do-site>

Um push com um commit sem ela (ou de outro site, ou de um usuário que o site
não conhece) é recusado, e o git mostra a linha que falta. O **hook**
`commit-msg` acrescenta a linha a cada commit; em **Pelo git**, na página, está
o comando que o instala no seu clone:

    curl -fsSL -u o-seu-login …/site-files/commit-msg -o .git/hooks/commit-msg && chmod +x .git/hooks/commit-msg

Para assinar commits já feitos: `git rebase -x 'git commit --amend --no-edit'
<o commit anterior a eles>`.

## O rascunho de outro usuário

Um rascunho compartilhado com você (passo **Trabalhar em conjunto**) também
é clonado e enviado pelo git: `…/site-files-drafts/<a chave do dono>.git` — o
endereço está na página do rascunho. Os seus commits continuam seus.

## Clonar

Clonar entrega **todos** os arquivos: por isso pede `@get` de cada um. Se um
caminho não pode ser visto (uma negação de `@get` nele), o clone é recusado.

## Com o desenvolvimento local

O repositório do site é o mesmo do deploy (passo **Sincronizar com o
desenvolvimento local**): um push aqui é como uma publicação — quem desenvolve
traz o que foi publicado antes do próximo deploy.
