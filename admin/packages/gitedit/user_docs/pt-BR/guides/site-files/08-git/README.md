# Pelo git

Os mesmos arquivos podem ser clonados e enviados (push) com o **git**, do seu
computador, com o seu **login** e uma **chave de acesso** (veja abaixo) — ou a
sua senha do admin. Há dois repositórios — as
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

## Com uma chave de acesso

Em vez da sua senha, use uma **chave de acesso** — o jeito para o git: ela
funciona mesmo quando você entra no admin pelo Google ou com o segundo fator,
pode ser só para os arquivos do site, expira e é revogada sem mexer na sua
senha. Crie a sua em {%= admin.model("my_access_keys").link %}, com as
permissões dos arquivos do site (`admin:/site-files:@get`, `!git` e as ações
dos arquivos que você vai mudar); o **código** aparece uma única vez —
copie-o.

### Clonar

O git pede o usuário e a senha: o usuário é o seu **login**; a **senha**, o
**código** da chave.

```
git clone https://<o-site>/admin/site-files-draft.git
Username for 'https://<o-site>': <o-seu-login>
Password for 'https://<o-seu-login>@<o-site>': <o-código-da-chave>
```

### Não pedir de novo

Sem mais nada, o git pede o código a cada `pull` e `push`. Guarde-o no
**gerenciador de credenciais** do git, uma vez para o site:

- **No arquivo** `~/.git-credentials` (texto, só legível por você):

  ```
  git config --global credential.https://<o-site>.helper store
  ```

  O próximo `git pull` (ou `push`) pede o usuário e o código e os grava; os
  seguintes não pedem mais. O arquivo fica com uma linha
  `https://<o-seu-login>:<o-código>@<o-site>`.

- **No chaveiro do sistema**, cifrado: no lugar de `store`, `osxkeychain`
  (macOS), `manager` (Windows, o Git Credential Manager) ou `libsecret`
  (Linux, quando instalado).

Não ponha o código **na URL** (`https://login:código@…`): ele fica gravado
no `.git/config` do clone, à vista, e vai junto se o clone for copiado.

### Quando a chave expira (ou é revogada)

O git passa a responder `Authentication failed`. Crie uma **chave nova** em
{%= admin.model("my_access_keys").link %} e troque a guardada:

1. Esqueça a antiga:

   ```
   printf 'protocol=https\nhost=<o-site>\n\n' | git credential reject
   ```

   (com `store`, também dá para editar a linha do site em
   `~/.git-credentials`.)

2. O próximo `git pull` pede o usuário e o código: informe o **novo** — ele é
   guardado no lugar do antigo.

Se o código estava na URL do clone, tire-o de lá:

```
git remote set-url origin https://<o-site>/admin/site-files-draft.git
```

Cada uso de uma chave fica no histórico dela e no seu registro de acessos.

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
