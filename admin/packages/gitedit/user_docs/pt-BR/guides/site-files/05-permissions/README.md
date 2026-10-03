# Permissões

Cada ação da página é uma permissão dela, dada aos papéis em **Papéis de
Usuário** (veja as políticas e permissões):

| Ação | Permissão |
| --- | --- |
| Ver a página e os arquivos | `@get` |
| Mudar os arquivos do rascunho | `!edit` |
| Commit | `!commit` |
| Atualizar | `!update` |
| Publicar | `!publish` |
| Descartar | `!discard` |
| Recomeçar | `!reset` |
| Pré-visualizar | `!preview` |

Assim, publicar pode ficar com menos gente que editar: quem edita faz o
commit, e outra pessoa, que revisa, publica.
