# Sistema de arquivos

Os arquivos do servidor que o site guarda: navegue pelas pastas, envie
arquivos e pastas (arraste-os para a lista), crie, edite, renomeie, copie e
apague; baixe um arquivo, ou uma pasta como zip; extraia um arquivo
compactado (zip, tar, gz).

O endereço acompanha a pasta aberta (`?dir=/pasta`): um link abre a página
nela.

## WebDAV

Os mesmos arquivos podem ser acessados de um computador por **WebDAV**, com o
endereço e a conta mostrados na página — como uma unidade de rede: o último
ícone da barra do explorador abre o painel **Protocolo WebDAV**, à direita
(expansível; a largura se arrasta pela borda). Lá estão o endereço, a conta
(o seu e-mail do admin; a senha, de preferência uma chave de acesso do seu
perfil, que pode ser revogada), as permissões que a chave precisa, as
pastas — e quais são só de leitura — e programas para usar.

### A chave de acesso

As permissões de uma chave só **negam**, sobre o que você pode: uma chave sem
nenhuma linha alcança o que você alcança. Para o WebDAV, as linhas dela não
podem negar o recurso desta página (o painel mostra o nome dele) em `@get`
(ler e baixar) nem, para escrever, em `!write` (enviar, mudar, mover e
apagar). Uma chave que só lê:

```
deny admin:/fs-tools:!write
```

## Terminal

O botão de terminal da barra (ou **F4**) abre um shell abaixo da lista, na
pasta aberta, e o acompanha quando se muda de pasta. A barra entre a lista e
o terminal se arrasta para mudar a altura dele (duplo clique: a altura
padrão; ↑ ↓ pelo teclado), guardada no navegador; o botão de expandir do
terminal o põe na janela toda, e o devolve.

Ele roda isolado do resto do servidor:

- vê só as pastas do sistema de arquivos, em `$FILES_ROOT`
  (`cd "$FILES_ROOT"/media`), cada uma lida ou escrita como você pode;
- o seu `HOME` e o `/tmp` são uma pasta da sessão, apagada ao fechar;
- não vê a configuração nem os segredos do servidor, nem a rede (a não ser
  que a configuração permita);
- tem limites de tempo de CPU, memória e tamanho de arquivo;
- fecha sozinho depois de 30 minutos sem uso, e ao fechar a página.

## Permissões

As permissões desta página valem para o explorador, o WebDAV e o terminal:

| Permissão | O que permite |
| --- | --- |
| `@get` | ver e baixar os arquivos |
| `!write` | enviar, criar, mudar, mover e apagar |
| `!terminal` | abrir o terminal — não está em nenhum conjunto de permissões: dê-a só a quem precisa |

Ao lado das pastas do servidor, aparecem pastas de cada usuário: o rascunho
dos **Arquivos do site** (`site-files`), com as permissões de editá-lo no
admin. Uma pasta configurada só de leitura nunca é escrita.

## Para o administrador do servidor

As pastas vêm de `RVQ_FSTOOLS_DIRS`, no `env:` do `app.yml`: `;` entre
elas, cada uma `MODO,PONTO,PASTA` — `r` só leitura, `w` leitura e escrita:

```yaml
env:
  RVQ_FSTOOLS_DIRS: "w,/media,data/media;r,/logs,/var/log/app"
```

O terminal é uma funcionalidade, `rvq.fstools.terminal`, **ligada por
padrão**. Para desligá-lo, no `app.yml` (e reinicie):

```yaml
features:
  rvq.fstools.terminal: false
```

Um programa compilado com a tag `rvq.fstools.noterminal` não o tem.

O isolamento e os limites, em `fstools.terminal` do `app.yml` (todas as
chaves opcionais; aqui, os padrões):

```yaml
fstools:
  terminal:
    sandbox: auto    # auto, bwrap, landlock ou off (sem isolamento)
    network: false   # o shell alcança a rede
    shell: ""        # padrão: /bin/bash, senão /bin/sh
    limits:
      cpu: 10m
      memory: 2G
      file_size: 1G
      processes: 0   # conta todos os processos do usuário do servidor
    idle: 30m
```

`auto` usa o bubblewrap, se ele roda no servidor, senão o Landlock do kernel
Linux; sem nenhum, o terminal fica desligado e o painel do WebDAV diz por
quê. Com o Landlock antes do Linux 6.12, o shell ainda pode mandar sinais aos
processos do usuário do servidor: rode o servidor com um usuário só dele, ou
desligue o terminal se isso importa.
