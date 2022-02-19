
**Utilizar o arquivo CSV que esta em gounico\csv, que contém uma pequena higienização.**
**O arquivo original esta faltando uma "," ao final da ultima linha. Essa foi a unica higienização manual**


**1 - Instalar Imagem / Subir Docker do MYSQL**

>   sudo docker pull mysql/mysql-server:latest

>     sudo docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root -d mysql/mysql-server:latest

**2 - Incluir usuario root para acesso local**
>     sudo docker exec -it mysql bash
>
>     mysql -uroot -p


    CREATE USER 'root'@'%' IDENTIFIED BY 'root';  
    FLUSH PRIVILEGES;  

**3 - Conferir e ver se está registrado**


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

**6 - Executar na pasta do projeto**

    go run gounico

**6 - Usar a seguinte collection do postman para testar as requisições**  
https://www.getpostman.com/collections/957a447bf220f3a9fab1