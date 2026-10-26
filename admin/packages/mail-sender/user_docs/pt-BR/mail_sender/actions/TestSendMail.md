# Testar envio de e-mail

Envia um e-mail de teste, pelo método escolhido no formulário da ação (o
Gmail ou o SMTP configurado), para o endereço informado, com o assunto e a
mensagem dados: um jeito rápido de conferir a conta antes de o site depender
dela. Um erro diz por que o e-mail não pôde ser enviado (uma senha errada,
uma porta bloqueada…).

O e-mail de teste não conta para os limites por endereço. Com
`RVQ_MAIL_SMTP` definida (desenvolvimento), ele vai para o servidor de
testes, não para o endereço de verdade.
