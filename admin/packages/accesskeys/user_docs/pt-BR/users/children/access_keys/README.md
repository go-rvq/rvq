# Chaves de acesso

As chaves de acesso deste usuário: códigos com que automações — git, WebDAV,
scripts — entram em nome dele, cada uma podendo o que o usuário pode, menos o
que ela nega. Quem pode editar o usuário cuida delas aqui; o próprio
usuário, em {%= admin.model("my_access_keys").link %}, onde tudo é explicado.

As permissões de uma chave são negações — `deny admin:users:*`, uma por
linha —, por cima do que o usuário pode; explicadas lá.
