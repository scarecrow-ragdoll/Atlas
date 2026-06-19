flowchart LR
  SUPPLIER["Поставщик\nRedeem.io\n(G-engine)"]

  BALANCE["Зачисление средств\nв ЛК поставщика"]

  DIGITAL["Выдача Digital Code\n(Steam и др.)"]

  SUPPLIER -->|"Получает перевод\nот Redeem.io (Crypto)"| BALANCE

  BALANCE -->|"Ручная покупка\nгифтов Redeem.io'ом"| DIGITAL

  BALANCE -->|"API покупка"| DIGITAL

  DIGITAL -->|"API доставка"| SUPPLIER

  DIGITAL -->|"Ручная доставка\nRedeem.io сам забирает"| SUPPLIER
