flowchart LR
  SELLER["Продавец\nRedeem.io"]

  PURCHASE["Redeem.io закупка\nDigital Code / AAAcode"]

  MP_SELL["Digital Marketplace\ng2a · eneba · etc\n— — —\nУспешная продажа покупателю"]

  MP_DELIVER["Digital Marketplace\n— — —\nВыдача Redeem.io\nваучера покупателю"]

  SELLER -->|"Пополнение баланса\nRedeem.io (Crypto)"| PURCHASE

  PURCHASE -->|"Запуск товара\nВыставление на продажу"| MP_SELL

  MP_SELL -->|"Ручная выдача\n(предварительная\nзагрузка товара)"| MP_DELIVER

  MP_SELL -->|"API выдача"| MP_DELIVER

  MP_DELIVER -->|"Вывод средств\nс Digital Marketplace\n(Crypto)"| SELLER

  MP_SELL -->|"Маржа 2%"| SELLER

  style MP_SELL stroke:#333
