# Link de download (.tar.gz)

Cria um link que baixa o **state inteiro** num `.tar.gz`, **sem login**, por
**24 horas**. Serve para levar o backup a outro lugar sem montar um
cliente: um navegador, `curl` ou `wget` baixam o arquivo direto.

O link aparece **só uma vez**, no diálogo, com os botões **Copiar** e
**Abrir** (abre no navegador e começa o download). Depois, a aplicação
guarda só uma soma dele: ninguém o recupera. Para outro, crie outro.

- Ele pode ser usado **mais de uma vez** dentro das 24 horas. O `.tar.gz` é
  gerado enquanto é enviado e não retoma de onde parou: se o download cair,
  comece de novo pelo mesmo link.
- Cada uso fica registrado em **Links** (quantos, o último endereço e o
  programa), onde o link também pode ser **revogado**.
- O link para de funcionar quando expira, quando é revogado ou quando o
  state é excluído.
- O arquivo mantém permissões, datas e links simbólicos. Um arquivo
  alterado no servidor depois do state fica de fora, e o `.tar.gz` termina
  com `BACKUP-SKIPPED.txt`, a lista deles.

Quem tem o link tem **todos os dados e segredos da aplicação**. Mande-o só
por um canal seguro e revogue-o quando não precisar mais.

```
curl -fLo backup.tar.gz 'https://seu-site/admin/!backup/<id>/archive/<token>.tar.gz'
tar -xzf backup.tar.gz
```

## Quando está disponível

Com o state **pronto**, para quem tem a permissão **Baixar** desta página.
