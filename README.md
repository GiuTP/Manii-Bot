# Manii-Bot

Um bot de Telegram desenvolvido em Go para controle financeiro pessoal e gestão de faturas de cartão de crédito

O Manii-Bot foi criado para substituir planilhas complexas por uma interface conversacional rápida. Ele recebe comandos simples pelo dispositivo, calcula projeções de parcelamentos, gerencia dias de fechamentos de faturas e consolida gastos mensais em relatórios detalhados.

Projeto está pronto para uso, apesar de está em desenvolvimento, então bugs ou erros podem ser encontrados.

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
* GCC instaldo (necessário para o driver do SQLite em Go)

### Instalação
1. Clone o repositório:
```bash
git clone git@github.com:GiuTP/Manii-Bot.git
cd manii-bot
```

2. Baixe as dependências:
```bash
go mod tidy
```

3. Configure as variáveis de ambiente (crie um arquivo `.env` na raiz):
```.env
TELEGRAM_TOKEN=SEU_TOKEN_AQUI
ALLOWED_USERS=ID_TELEGRAM
DB_PATH = caminho/para/seu/banco/de/dados
```

4. Execute o bot:
```bash
go run cmd/bot/main.go
```

*(Na primeira execução, o bot criará automaticamente o arquivo finance.db e as tabelas necessárias)*

## :book: Como usar

Envie comandos diretamente para o seu bot no Telegram. Exemplos:
* **Criar cartão:**`/cartao banco credito 25 2` (nome, tipo, fechamento, vencimento)
* **Registrar compra:**`/compra [LOJA] Compra A 3x 5.00 banco pessoa 01/01/2007`
* **Listar fatura:**`listar e 8` (listar todas as movimentações do mês 8)
* **Extrato de uma pessoa:** `listar e p pessoa 8`

Para a listar completa junto de suas sintaxe, mande `/help` para o bot.

# :page_facing_up: Licença
Este projeto está sob a licença [ GPLv3](LICENSE).