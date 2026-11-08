# Grupos de usuários

Um grupo reúne usuários que trabalham do mesmo jeito: dê ao grupo papéis, ou
compartilhe um registro com ele, uma vez só — e todo membro tem.

- **Membros** — os usuários do grupo. Quem entra alcança na hora o que o
  grupo alcança; quem sai perde na hora.
- **Papéis** — os papéis do grupo, **em ordem** (arraste para mudar), como
  os de um usuário (veja {%= admin.model("users").link %}): nas permissões
  de um membro, eles vêm depois dos papéis do próprio membro, os grupos na
  ordem dos grupos do membro.

Um registro compartilhado com um grupo (uma organização, por exemplo) não
precisa de convite: o proprietário confirma por e-mail, e daí em diante
todo membro o alcança. O que é compartilhado com o próprio membro vale antes
do que é compartilhado com os grupos dele.

Excluir um grupo tira os membros dele — cada um perdendo o que tinha por ele.
