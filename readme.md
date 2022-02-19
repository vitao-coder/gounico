
>**Requisitos: Docker ou Docker For Windows com WSL2 e/ou Go 1.17**

>**Utilizar o arquivo CSV que esta em gounico\csv, que contém uma pequena higienização.**
**O arquivo original esta faltando uma "," ao final da ultima linha. Essa foi a unica higienização manual**

**1 - Instalar Imagem / Subir Docker do MYSQL**

>   sudo docker pull mysql/mysql-server:latest

>   sudo docker run -p 3306:3306 --add-host host.docker.internal:host-gateway --name mysql -e MYSQL_ROOT_PASSWORD=root -d mysql/mysql-server:latest

**2 - Incluir usuario root para acesso local**
>     sudo docker exec -it mysql bash

>     mysql -uroot -p

>  CREATE USER 'root'@'%' IDENTIFIED BY 'root';

>  FLUSH PRIVILEGES;  

**3 - Conferir se está conforme abaixo**


    SELECT user, host FROM mysql.user;  

    +------------------+-----------+  
    | user             | host      |  
    +------------------+-----------+  
    | root             | %         |  
    | healthchecker    | localhost |  
    | mysql.infoschema | localhost |  
    | mysql.session    | localhost |  
    | mysql.sys        | localhost |  
    | root             | localhost |  
    +------------------+-----------+  


**4 - Criar schema da aplicação**


    CREATE SCHEMA `gounico` ;  


**5 - Dar permissao full ao usuario root no esquema**

    GRANT ALL PRIVILEGES ON gounico.* TO 'root'@'%' WITH GRANT OPTION;

**6 - Caso deseje executar utilizando o próprio GO (se tiver instalado na maquina),  executar o comando abaixo na pasta do projeto**

    go run gounico

 **6.1 - Caso deseje executar utilizando o docker utilizar os comandos abaixo, para criar a imagem e executar o container**
> docker build . -t gounico
> 
> sudo docker run -itd -p 8000:8000 --add-host host.docker.internal:host-gateway --name gounicoApp gounico
> 
**7 - Usar a seguinte collection do postman para testar as requisições**  
> https://www.getpostman.com/collections/957a447bf220f3a9fab1
>
> Endpoints implementados:
> 
> > POST http://localhost:8000/csvprocessor - Carregamento do arquivo CSV e normalização dos dados
> 
> > POST http://localhost:8000/novafeira - Criação de noova feira
> 
> > GET http://localhost:8000/buscarfeira/bairro/JD BRASILIA - Buscar feiras por bairro
> 
> > GET http://localhost:8000/buscarfeira/distrito/ARICA - Buscar feiras por distrito
> 
> > DELETE http://localhost:8000/excluirfeira/98 - Deletar uma feira por id
> 
> > ~~PUT http://localhost:8000/alterarfeira/1~~ - **NOT IMPLEMENTED**

> **Observações**: Não foi implementado o PUT para Feira (Em WIP) e os testes unitários estão incompletos (WIP).