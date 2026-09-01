# Manii-Bot

Um bot de Telegram desenvolvido em Go para controle financeiro pessoal e gestão de faturas de cartão de crédito

O Manii-Bot foi criado para substituir planilhas complexas por uma interface conversacional rápida. Ele recebe comandos simples pelo dispositivo, calcula projeções de parcelamentos, gerencia dias de fechamentos de faturas e consolida gastos mensais em relatórios detalhados.

O projeto está pronto para uso, apesar de estar em desenvolvimento, então bugs ou erros podem ser encontrados.

## :rocket: Funcionalidades

* **Cadastro rápido de despesas:** registre compras avulsas ou parceladas com uma única mensagem.
* **Cálculo inteligente de fatura:** o bot entende dias de fechamento e vencimento de diferentes cartões, jogando parcelas automaticamente para os meses corretos.
* **Gestão de compras recorrentes:** controle de gastos fixos recorrentes, por exemplo streaming.
* **Extratos mensais:** geração de relatórios filtrados por mês, pessoa ou cartão.
* **Multi-usuário:** suporte para dividir gastos entre diferentes pessoas.

## :hammer: Tecnologias Utilizadas
* **Linguagem**: Go
* **Banco de Dados:** SQLite3
* **Arquitetura**: baseada em DDD
* **Dependências** bibliotecas nativas, `go-sqlite3`, `gofpdf`, `godotenv`, `go-telegram-bot-api`.

## :package: Como Rodar Localmente

### Pré-requisitos
* Go 1.21 ou superior instalado
* Token de um bot no Telegram
* GCC instalado (necessário para o driver do SQLite em Go)

### Instalação
1. Clone o repositório:
```bash
git clone git@github.com:GiuTP/Manii-Bot.git
cd Manii-Bot
```

2. Baixe as dependências:
```bash
go mod tidy
```

3. Configure as variáveis de ambiente (crie um arquivo `.env` na raiz):
```.env
TELEGRAM_TOKEN=SEU_TOKEN_AQUI (via @BotFather)
ALLOWED_USERS=ID_TELEGRAM (via @userinfobot)
DB_PATH=caminho/para/seu/banco/de/dados
```

4. Execute o bot:
```bash
go run cmd/bot/main.go
```

*(Na primeira execução, o bot criará automaticamente o arquivo finance.db e as tabelas necessárias)*

## :book: Como usar

Envie comandos diretamente para o seu bot no Telegram. Exemplos:
* **Criar cartão:**`/cartao banco credito 25 2` (nome, tipo, fechamento, vencimento)
* **Registrar compra:**`/compra Compra A 3x 5.00 banco pessoa 01/01/2007`
* **Listar fatura:**`/listar e 8` (listar todas as movimentações do mês 8)
* **Extrato de uma pessoa:** `listar e p pessoa 8`

Para a lista completa dos comandos junto de suas sintaxes e detalhamentos, mande `/help` para o bot.

### :mag: Demonstração
Um exemplo visual de como é a conversa com o bot. Foram feitos os seguintes comandos:

1. **Criação de entidades:** criação da pessoa (`JoãoM`) e do cartão (`CartãoT`).
2. **Registro de despesa:** inserção de uma compra parcelada vinculando a pessoa e o cartão
3. **Consulta e filtros:** listagem filtrada por responsável (`JoãoM`).
4. **Relatório**: geração do extrato mensal em PDF.
5. **Manutenção:** atualização da descrição da despesa e posteriomente remorção do registro.

![Interação com o bot no Telegram](docs/screenshots/tela_bot.png)

### Layout do Relatório Gerado
O extrato exportado organiza os lançamentos com ordenação, detalhes de parcelamento e linhas divisórias:

![Extrato em PDF](docs/screenshots/tela_extrato.png)

# :page_facing_up: Licença
Este projeto está sob a licença [GPLv3](LICENSE).
