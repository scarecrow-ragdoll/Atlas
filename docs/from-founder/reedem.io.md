flowchart LR
  RIO["Redeem.io"]

  BALANCE["Redeem.io\nзачисление баланса\nв ЛК продавца"]

  FORMING["Redeem.io\nформирование\nDigital Code\n(собственный код)"]

  SELLER["Продавец\nRedeem.io"]

  SELLER -->|"Пополнение баланса\n(Crypto)"| RIO

  RIO --> BALANCE

  BALANCE -->|"Ручная покупка\nвыбранных номиналов"| FORMING

  BALANCE -->|"API запрос на покупку\nвыбранных номиналов"| FORMING

  FORMING -->|"Ручное получение"| SELLER

  FORMING -->|"API получение"| SELLER

  SELLER -->|"Продавец успешно продал\nна Digital Marketplace\nи получил деньги обратно"| RIO
