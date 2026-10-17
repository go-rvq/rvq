# Passo a passo

Como editar os arquivos de que o site é feito — templates, estilos, scripts,
imagens — pelo admin, em **{%= admin.page("/site-files").menu %}**, sem
mexer no site até a hora de publicar; e como isso convive com quem
desenvolve o site no próprio computador.

O caminho, passo a passo:

1. **Editar no rascunho** — o seu, só seu.
2. **Pré-visualizar** — o site como o rascunho o faz.
3. **Commit** — registrar as alterações, com uma mensagem.
4. **Publicar** — pôr os commits no site.
5. **Permissões** — quem pode cada coisa.
6. **Sincronizar com o desenvolvimento local** — os commits feitos aqui e os
   feitos no computador de quem desenvolve.
7. **Opções dos layouts** — os campos dos formulários delas, e um select de
   valores com `options`.
8. **Pelo git** — clone e faça push do rascunho, ou do site, do seu computador.
9. **Trabalhar em conjunto** — compartilhe o seu rascunho.
10. **Os arquivos padrão** — comece um repositório vazio com eles, ou
    restaure-os.

**Na primeira vez**, sem repositório em `public/`, a aplicação cria um quando
sobe. Se `public/` não existe ou está vazio, o repositório fica vazio, e a
página mostra como dar a ele os primeiros arquivos: os padrão, ou um push
pelo git (veja **Os arquivos padrão**). Se `public/` já tem arquivos, eles
entram num primeiro commit, `first commit`, sem os ocultos (`.gitignore` com
`.*`). Quem decide é o administrador do servidor (`PUBLIC_GIT_INIT`).
