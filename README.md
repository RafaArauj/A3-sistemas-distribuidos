# Sistema Distribuído — Processamento de Texto

Aplicação cliente-servidor em Go em que o **cliente** recebe um texto do usuário e distribui o processamento entre **dois servidores**, que executam operações diferentes em paralelo. A comunicação é feita via **socket TCP** com mensagens em **JSON**.
