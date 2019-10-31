# GO ImageServer

## Inicialização

 execute ```cp variables.example.env variables.env``` e preencha o arquivo com os dados para a conexão com o banco de dados

 execute ```docker-compose build && docker-compose up``` para criar as imagens e subir os containeres
    
## Rotas

### /upload :
-  POST:
    - Faz o upload de uma imagem, gerando alguns arquivos em tamanhos pré-definidos

### /imagens :
- GET:
    - Lista todas as imagens cadastradas
    
            
### /imagen-info/{uuid} :
- GET:
    - Retorna um json com os dados da imagem original
     ```json
        {
           "ID":98, //id na tabela arquivo
           "ImagemID":7, //referência à imagem original
           "Tamanho":"800x800", //tamanho do arquivo (ALTURA x LARGURA)
           "Path":"temp-images/img-c41ed461-39ed-65a9-85cc-724551ae3c97-800x800.jpeg", //caminho onde a imagem está salva
           "Original":false //true se o arquivo for o original
        }
     ```

### /imagem-info/{uuid}/{tamanho} :
- GET:
    - Retorna um json com os seguintes dados :
     ```json
        {
           "ID":98, //id na tabela arquivo
           "ImagemID":7, //referência à imagem original
           "Tamanho":"800x800", //tamanho do arquivo (ALTURA x LARGURA)
           "Path":"temp-images/img-c41ed461-39ed-65a9-85cc-724551ae3c97-800x800.jpeg", //caminho onde a imagem está salva
           "Original":false //true se o arquivo for o original
        }
     ```

### /imagem/{uuid} :
- GET :

    Retorna o arquivo original referente ao parâmetro uuid 

### /imagem/{uuid}/{tamanho} :
- GET :

    Retorna o arquivo referente ao parâmetro uuid com o tamanho solicitado, se o arquivo não existir nesse tamanho, o mesmo será gerado e retornado 