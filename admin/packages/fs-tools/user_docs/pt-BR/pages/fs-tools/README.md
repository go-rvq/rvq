# Sistema de arquivos

Os arquivos do servidor que o site guarda: navegue pelas pastas, envie
arquivos e pastas (arraste-os para a lista), crie, edite, renomeie, copie e
apague; baixe um arquivo, ou uma pasta como zip; extraia um arquivo
compactado.

Os mesmos arquivos podem ser acessados de um computador por **WebDAV**, com o
endereço e a conta mostrados na página — como uma unidade de rede: no painel
**Protocolo WebDAV** acima do explorador, ou no botão do fim da sua barra, que
diz também quais pastas são só de leitura.

## Permissões

O WebDAV pede as permissões desta página: ver e baixar os arquivos pedem
`@get`; enviar, mudar, mover e apagar, `!write`.

Ao lado das pastas do servidor, o WebDAV mostra pastas de cada usuário: o
rascunho dos **Arquivos do site** (`site-files`), com as permissões de
editá-lo no admin.

## Terminal

O botão de terminal da barra (ou **F4**) abre um shell abaixo da lista, na
pasta aberta, e o acompanha quando se muda de pasta. Ele vê só as pastas do
sistema de arquivos — em `$FILES_ROOT` —, cada uma lida ou escrita como o
usuário pode; nada mais do servidor, nem a rede (a não ser que a
configuração permita).

Pede a permissão `!terminal` desta página, que não está em nenhum conjunto de
permissões: dê-a a quem precisa. O servidor pode desligá-lo — o painel do
WebDAV diz por quê a quem tem a permissão.
