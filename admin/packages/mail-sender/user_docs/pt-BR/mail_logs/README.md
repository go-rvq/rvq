# Central de e-mails

Todo e-mail que o site envia, ou tenta enviar, fica registrado aqui: as
mensagens dos formulários de contato, as confirmações, os avisos de
comentários, as redefinições de senha, o e-mail de teste. É só para
consultar: ninguém edita, cria ou apaga um registro (os antigos saem pela
retenção, abaixo).

**A listagem** mostra o mais novo primeiro, com a situação (**Enviado** ou
**Falhou**), o tipo, para quem, o assunto e o método. A aba **Falhas** reúne
os que não saíram. Os filtros escolhem por tipo, período, destinatário e
situação; a busca acha pelo assunto, pelo destinatário, pela referência do
registro de origem e pelo Message-ID.

**O detalhe** de cada e-mail mostra:

- quando, a situação e, se falhou, **por quê** (uma senha errada, o limite
  por endereço, o envio não configurado…);
- **Origem**: o botão **Abrir origem** leva ao registro de que o e-mail é —
  a mensagem de contato, a resposta, o usuário;
- para quem, responder para, o assunto e o método (Gmail, SMTP, ou o servidor
  de testes no desenvolvimento);
- o Message-ID e a mensagem a que responde: a **conversa**. Abaixo do corpo,
  os outros e-mails da mesma conversa;
- quem provocou o envio, de que IP e navegador, e o idioma;
- o **corpo**, como foi enviado, isolado da página.

## O que nunca fica guardado

O conteúdo sensível é ocultado antes de gravar e aparece como `••••••`: os
valores que quem envia marca como sensíveis (um link de redefinição de
senha, um telefone) e, em qualquer link, o valor de parâmetros como `token`,
`key`, `code` e `password`. Alguns e-mails nem guardam o corpo
(**[conteúdo oculto]**). Um e-mail com algo ocultado não pode ser enviado de
novo: o que ele tinha não está mais aqui.

## Por quanto tempo

A configuração do {%= admin.model("mail_sender").link %}, aba **Envio**, diz
por quantos dias cada e-mail fica guardado (padrão 365; `0`, para sempre).
Os mais antigos são apagados uma vez por dia.
