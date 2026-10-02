# Detalhe: {%= doc.model %}

O que um registro mostra quando é aberto pela listagem.

{%= admin.fields("detail") %}

## O menu do registro

O menu (**⋯**) no topo do detalhe: o que está dentro do registro, o que pode
ser feito com ele, e as páginas dele. Cada item aparece só para quem pode
usá-lo; um disponível conforme o registro diz quando.

{%= admin.menu() %}
