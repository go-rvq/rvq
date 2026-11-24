# Tamanhos das imagens

Cada campo de imagens do admin — a capa de uma página, a imagem de um
slider, o logo do site… — tem os seus **tamanhos**: os recortes que o site
mostra da imagem escolhida ali, cada um com um nome (o que os templates do
site pedem), uma largura e uma altura.

![Os tamanhos das imagens](images/listing.png)

- Os tamanhos começam como os **do site** (escritos no código dele). Edite um
  campo para definir os **seus**: a tabela de tamanhos — nome, largura,
  altura —; a lista diz quando os tamanhos em uso são os seus.
- Salvar tamanhos novos **refaz os recortes** de todos os arquivos que o
  campo mostra, em segundo plano: a coluna **Recortes** diz até onde foi (em
  andamento, prontos, ou o que deu errado). A área de um recorte é mantida
  quando a proporção é a mesma; senão, usa-se a maior área da proporção nova,
  centralizada na imagem — recorte de novo na biblioteca de mídia para
  escolher outra.
- Dois campos só podem ter um tamanho de mesmo nome com as **mesmas
  medidas**: senão, um arquivo usado nos dois seria recortado em um, depois
  no outro.
- Para voltar aos tamanhos do site, deixe a tabela como a do site; quando os
  tamanhos do próprio site mudam numa versão nova, os recortes são refeitos
  sozinhos.

## Formatos

Um campo pode oferecer **formatos**: os seus tamanhos em outra proporção — a
capa de uma página em retrato, quando as outras são em paisagem. Cada formato
tem um nome (letras minúsculas, dígitos, `-`), um rótulo e o seu maior
tamanho; os outros tamanhos seguem a proporção dele. Um formato `retrato` de
1080×1440 num campo de 1920×1080, 900×506 e 500×281 recorta em 1080×1440,
506×675 e 281×375.

- No registro, a imagem do campo mostra **Formato**: **Padrão** (os tamanhos
  do campo) ou um dos formatos. Escolher um refaz os recortes da imagem no
  formato, na hora.
- O site pede os mesmos tamanhos de sempre (`site.full`): o registro de um
  formato recebe os recortes dele. Nada muda nos templates.
- Salvar formatos novos refaz os recortes dos arquivos do campo também nos
  formatos.
