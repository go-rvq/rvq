# Chaves de acesso

As chaves de acesso deste usuário: códigos com que automações — git, WebDAV,
scripts — entram em nome dele, cada uma podendo só o que libera do que o
usuário pode. Quem pode editar o usuário cuida delas aqui; o próprio
usuário, em {%= admin.model("my_access_keys").link %}, onde tudo é explicado.

As permissões de uma chave liberam o que o usuário pode — `admin:posts:*`,
`*` (tudo) —, negam uma parte (`deny admin:posts:*:@delete`) ou a devolvem
(`allow`); sem nenhuma, a chave não alcança nada. São editadas no campo de
permissões, montado do que o usuário pode; explicadas lá.
