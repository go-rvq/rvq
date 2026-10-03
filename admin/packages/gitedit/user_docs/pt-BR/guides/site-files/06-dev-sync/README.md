# Sincronizar com o desenvolvimento local

Os arquivos do site são um repositório **git**. Ele está em dois lugares: no
**servidor** — o que o site usa e o que este editor muda — e no computador de
quem **desenvolve** o site (o `public/` do projeto). Os commits feitos aqui
ficam no servidor; os feitos no computador chegam ao servidor no deploy. Para
os dois lados não se perderem, cada um traz o do outro antes de mandar o seu.

## O que foi feito no admin, para o computador

Antes de mexer no `public/` e, sobretudo, antes de um deploy, quem desenvolve
traz os commits do servidor:

```
pub repo public pull server main
```

(o `pub` roda o git do `public/` com o servidor como `server`). Depois,
registra no projeto o novo commit do `public/` (o submódulo):

```
git add public && git commit -m "public: o que foi publicado pelo admin"
```

Sem isso, o deploy é recusado: o servidor tem commits que o computador não
tem, e o deploy não sobrescreve o que foi publicado pelo admin.

## O que foi feito no computador, para o admin

O deploy (`pub run`) publica o commit do `public/` do projeto no servidor: o
site passa a usá-lo. Quem edita no admin traz esses commits para o seu
rascunho com **Atualizar**.

## Arquivos mudados no servidor

Se um arquivo do site for mudado direto no servidor, fora do git, tanto o
**Publicar** do admin quanto o deploy param e avisam, sem sobrescrever nada:
quem cuida do servidor decide o que fica (um commit com a mudança, ou
desfazê-la).
